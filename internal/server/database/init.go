package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

func Init() (*Queries, error) {
	dbFile, err := createDatabaseDirectory()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not init db: %w", err)
	}

	if _, err := db.ExecContext(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("could not execute schema: %w", err)
	}

	return New(db), nil
}

func createDatabaseDirectory() (string, error) {
	dbDir := os.Getenv("DB_DIR")
	if dbDir == "" {
		panic("DB_DIR env var is not set")
	}

	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return "", fmt.Errorf("could not create database dir: %w", err)
	}

	return filepath.Join(dbDir, "blkhell.db"), nil
}
