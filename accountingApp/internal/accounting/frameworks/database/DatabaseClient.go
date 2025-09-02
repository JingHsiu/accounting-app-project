package database

// DatabaseClient provides an abstraction layer for database operations
// This interface sits in Layer 2 (Application) and is implemented in Layer 4 (Frameworks)
// allowing Layer 3 (Adapter) implementations to use database operations without direct dependencies
type DatabaseClient interface {
	// QueryRow executes a query that is expected to return at most one row
	QueryRow(query string, args ...interface{}) RowScanner
	
	// Query executes a query that returns multiple rows
	Query(query string, args ...interface{}) (RowsScanner, error)
	
	// Exec executes a query that doesn't return rows (INSERT, UPDATE, DELETE)
	Exec(query string, args ...interface{}) (ExecResult, error)
	
	// BeginTx starts a new transaction
	BeginTx() (Transaction, error)
}

// RowScanner abstracts single row scanning operations
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// RowsScanner abstracts multiple rows scanning operations  
type RowsScanner interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

// ExecResult abstracts the result of an execution operation
type ExecResult interface {
	RowsAffected() (int64, error)
}

// Transaction provides transactional database operations
// It inherits all DatabaseClient operations and adds transaction-specific methods
type Transaction interface {
	DatabaseClient
	Commit() error
	Rollback() error
}