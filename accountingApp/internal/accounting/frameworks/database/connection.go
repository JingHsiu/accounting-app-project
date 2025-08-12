package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	db *sql.DB
}

func NewPostgresConnection(databaseURL string) (*PostgresConnection, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresConnection{db: db}, nil
}

func (pc *PostgresConnection) GetDB() *sql.DB {
	return pc.db
}

func (pc *PostgresConnection) Close() error {
	return pc.db.Close()
}