package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Init() (*Queries, error) {
	dbFile, err := createDatabaseDirectory()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not init db: %w", err)
	}

	return New(db), nil
}

// Open opens the database connection
func Open() (*sql.DB, error) {
	dbFile, err := createDatabaseDirectory()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	return db, nil
}

func createDatabaseDirectory() (string, error) {
	dbDir := os.Getenv("DB_DIR")
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return "", fmt.Errorf("could not create database dir: %w", err)
	}

	return filepath.Join(dbDir, "blkhell.db"), nil
}
