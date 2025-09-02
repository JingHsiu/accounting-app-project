package database

import (
	"database/sql"
)

// PostgreSQLClient implements the DatabaseClient interface for PostgreSQL
// This is a Layer 4 (Frameworks) implementation that wraps *sql.DB
type PostgreSQLClient struct {
	db *sql.DB
}

// NewPostgreSQLClient creates a new PostgreSQL database client
func NewPostgreSQLClient(db *sql.DB) DatabaseClient {
	return &PostgreSQLClient{
		db: db,
	}
}

// QueryRow executes a query that returns at most one row
func (c *PostgreSQLClient) QueryRow(query string, args ...interface{}) RowScanner {
	return &SQLRowWrapper{
		row: c.db.QueryRow(query, args...),
	}
}

// Query executes a query that returns multiple rows
func (c *PostgreSQLClient) Query(query string, args ...interface{}) (RowsScanner, error) {
	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLRowsWrapper{
		rows: rows,
	}, nil
}

// Exec executes a query that doesn't return rows
func (c *PostgreSQLClient) Exec(query string, args ...interface{}) (ExecResult, error) {
	result, err := c.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLExecResultWrapper{
		result: result,
	}, nil
}

// BeginTx starts a new transaction
func (c *PostgreSQLClient) BeginTx() (Transaction, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	return &PostgreSQLTransaction{
		tx: tx,
	}, nil
}

// SQLRowWrapper wraps *sql.Row to implement RowScanner interface
type SQLRowWrapper struct {
	row *sql.Row
}

func (w *SQLRowWrapper) Scan(dest ...interface{}) error {
	return w.row.Scan(dest...)
}

// SQLRowsWrapper wraps *sql.Rows to implement RowsScanner interface
type SQLRowsWrapper struct {
	rows *sql.Rows
}

func (w *SQLRowsWrapper) Next() bool {
	return w.rows.Next()
}

func (w *SQLRowsWrapper) Scan(dest ...interface{}) error {
	return w.rows.Scan(dest...)
}

func (w *SQLRowsWrapper) Close() error {
	return w.rows.Close()
}

// SQLExecResultWrapper wraps sql.Result to implement ExecResult interface
type SQLExecResultWrapper struct {
	result sql.Result
}

func (w *SQLExecResultWrapper) RowsAffected() (int64, error) {
	return w.result.RowsAffected()
}

// PostgreSQLTransaction implements Transaction interface for PostgreSQL
type PostgreSQLTransaction struct {
	tx *sql.Tx
}

// QueryRow executes a query within transaction that returns at most one row
func (t *PostgreSQLTransaction) QueryRow(query string, args ...interface{}) RowScanner {
	return &SQLRowWrapper{
		row: t.tx.QueryRow(query, args...),
	}
}

// Query executes a query within transaction that returns multiple rows
func (t *PostgreSQLTransaction) Query(query string, args ...interface{}) (RowsScanner, error) {
	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLRowsWrapper{
		rows: rows,
	}, nil
}

// Exec executes a query within transaction that doesn't return rows
func (t *PostgreSQLTransaction) Exec(query string, args ...interface{}) (ExecResult, error) {
	result, err := t.tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLExecResultWrapper{
		result: result,
	}, nil
}

// BeginTx is not supported within a transaction (nested transactions not supported by PostgreSQL driver)
func (t *PostgreSQLTransaction) BeginTx() (Transaction, error) {
	// PostgreSQL doesn't support nested transactions with the standard library
	// Return the current transaction instead
	return t, nil
}

// Commit commits the transaction
func (t *PostgreSQLTransaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *PostgreSQLTransaction) Rollback() error {
	return t.tx.Rollback()
}