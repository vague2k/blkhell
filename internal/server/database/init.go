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
	isDev := os.Getenv("GO_ENV") != "production"

	dbFile, err := resolveDBPath(isDev)
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

func resolveDBPath(isDev bool) (string, error) {
	if isDev {
		return createDatabaseDirectory(".")
	}

	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataHome = filepath.Join(home, ".local", "share")
	}

	return createDatabaseDirectory(dataHome)
}

func createDatabaseDirectory(base string) (string, error) {
	dir := filepath.Join(base, "blkhell")
	dbFile := filepath.Join(dir, "blkhell.db")

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("could not create database dir: %w", err)
	}

	return dbFile, nil
}
