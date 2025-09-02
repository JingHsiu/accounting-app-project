package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/JingHsiu/accountingApp/internal/accounting/adapter/store"
)

// PgAggregateStoreAdapter implements AggregateStore interface for PostgreSQL
// This is Layer 4 (Frameworks & Drivers) adapter that uses DatabaseClient
type PgAggregateStoreAdapter[T store.AggregateData] struct {
	dbClient  DatabaseClient
	tableName string
	columns   []string
	scanner   func(RowScanner) (*T, error)
	inserter  func(T) []interface{}
}

// NewPgAggregateStoreAdapter creates a new PostgreSQL aggregate store adapter
func NewPgAggregateStoreAdapter[T store.AggregateData](
	dbClient DatabaseClient,
	tableName string,
	columns []string,
	scanner func(RowScanner) (*T, error),
	inserter func(T) []interface{},
) store.AggregateStore[T] {
	return &PgAggregateStoreAdapter[T]{
		dbClient:  dbClient,
		tableName: tableName,
		columns:   columns,
		scanner:   scanner,
		inserter:  inserter,
	}
}

// Save persists aggregate state using upsert pattern
func (s *PgAggregateStoreAdapter[T]) Save(data T) error {
	if len(s.columns) == 0 {
		return fmt.Errorf("no columns defined for table %s", s.tableName)
	}

	// Build the INSERT ... ON CONFLICT DO UPDATE query
	placeholders := make([]string, len(s.columns))
	updateSet := make([]string, len(s.columns)-1) // Exclude ID from update

	for i, col := range s.columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		if i > 0 { // Skip ID column (first column)
			updateSet[i-1] = fmt.Sprintf("%s = EXCLUDED.%s", col, col)
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (%s)
		ON CONFLICT (id) DO UPDATE SET
			%s
	`, s.tableName,
		strings.Join(s.columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updateSet, ", "))

	values := s.inserter(data)
	_, err := s.dbClient.Exec(query, values...)
	return err
}

// FindByID retrieves aggregate state by ID
func (s *PgAggregateStoreAdapter[T]) FindByID(id string) (*T, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s 
		WHERE id = $1
	`, strings.Join(s.columns, ", "), s.tableName)

	row := s.dbClient.QueryRow(query, id)
	data, err := s.scanner(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	return data, nil
}

// Delete removes aggregate state by ID
func (s *PgAggregateStoreAdapter[T]) Delete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName)
	result, err := s.dbClient.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("aggregate with id %s not found in table %s", id, s.tableName)
	}

	return nil
}

// PgBatchAggregateStoreAdapter extends PgAggregateStoreAdapter with batch operations
type PgBatchAggregateStoreAdapter[T store.AggregateData] struct {
	*PgAggregateStoreAdapter[T]
}

// NewPgBatchAggregateStoreAdapter creates a new PostgreSQL batch aggregate store adapter
func NewPgBatchAggregateStoreAdapter[T store.AggregateData](
	dbClient DatabaseClient,
	tableName string,
	columns []string,
	scanner func(RowScanner) (*T, error),
	inserter func(T) []interface{},
) store.BatchAggregateStore[T] {
	baseStore := &PgAggregateStoreAdapter[T]{
		dbClient:  dbClient,
		tableName: tableName,
		columns:   columns,
		scanner:   scanner,
		inserter:  inserter,
	}

	return &PgBatchAggregateStoreAdapter[T]{
		PgAggregateStoreAdapter: baseStore,
	}
}

// SaveBatch persists multiple aggregates in a single transaction
func (s *PgBatchAggregateStoreAdapter[T]) SaveBatch(data []T) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := s.dbClient.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range data {
		if err := s.saveSingle(tx, item); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// saveSingle saves a single aggregate within a transaction
func (s *PgBatchAggregateStoreAdapter[T]) saveSingle(tx Transaction, data T) error {
	if len(s.columns) == 0 {
		return fmt.Errorf("no columns defined for table %s", s.tableName)
	}

	placeholders := make([]string, len(s.columns))
	updateSet := make([]string, len(s.columns)-1)

	for i, col := range s.columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		if i > 0 {
			updateSet[i-1] = fmt.Sprintf("%s = EXCLUDED.%s", col, col)
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (%s)
		ON CONFLICT (id) DO UPDATE SET
			%s
	`, s.tableName,
		strings.Join(s.columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updateSet, ", "))

	values := s.inserter(data)
	_, err := tx.Exec(query, values...)
	return err
}

// FindBatch retrieves multiple aggregates by their IDs
func (s *PgBatchAggregateStoreAdapter[T]) FindBatch(ids []string) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}

	// Build IN clause with placeholders
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s 
		WHERE id IN (%s)
	`, strings.Join(s.columns, ", "), s.tableName, strings.Join(placeholders, ", "))

	rows, err := s.dbClient.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		data, err := s.scanner(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil
}

// PgQueryAggregateStoreAdapter extends PgAggregateStoreAdapter with query capabilities
type PgQueryAggregateStoreAdapter[T store.AggregateData] struct {
	*PgAggregateStoreAdapter[T]
}

// NewPgQueryAggregateStoreAdapter creates a new PostgreSQL query aggregate store adapter
func NewPgQueryAggregateStoreAdapter[T store.AggregateData](
	dbClient DatabaseClient,
	tableName string,
	columns []string,
	scanner func(RowScanner) (*T, error),
	inserter func(T) []interface{},
) store.QueryAggregateStore[T] {
	baseStore := &PgAggregateStoreAdapter[T]{
		dbClient:  dbClient,
		tableName: tableName,
		columns:   columns,
		scanner:   scanner,
		inserter:  inserter,
	}

	return &PgQueryAggregateStoreAdapter[T]{
		PgAggregateStoreAdapter: baseStore,
	}
}

// FindBy retrieves aggregates matching given criteria
func (s *PgQueryAggregateStoreAdapter[T]) FindBy(criteria map[string]interface{}) ([]T, error) {
	if len(criteria) == 0 {
		return []T{}, nil
	}

	// Build WHERE clause from criteria
	var whereClauses []string
	var args []interface{}
	argIndex := 1

	for column, value := range criteria {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s 
		WHERE %s
	`, strings.Join(s.columns, ", "), s.tableName, whereClause)

	rows, err := s.dbClient.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		data, err := s.scanner(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil
}

// Count returns the number of aggregates matching criteria
func (s *PgQueryAggregateStoreAdapter[T]) Count(criteria map[string]interface{}) (int64, error) {
	if len(criteria) == 0 {
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", s.tableName)
		row := s.dbClient.QueryRow(query)
		var count int64
		err := row.Scan(&count)
		return count, err
	}

	// Build WHERE clause from criteria
	var whereClauses []string
	var args []interface{}
	argIndex := 1

	for column, value := range criteria {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", s.tableName, whereClause)

	row := s.dbClient.QueryRow(query, args...)
	var count int64
	err := row.Scan(&count)
	return count, err
}
