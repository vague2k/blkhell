package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

// MigrateUp applies all pending up migrations
func MigrateUp(db *sql.DB) error {
	m, err := migrator(db)
	if err != nil {
		return err
	}

	fmt.Println("Running migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Successfully applied migrations")
	return nil
}

// MigrateDown applies one down migration
func MigrateDown(db *sql.DB) error {
	m, err := migrator(db)
	if err != nil {
		return err
	}

	fmt.Println("Rolling back one migration...")
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Successfully rolled back migration")
	return nil
}

func migrator(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://server/database/migrations",
		"sqlite",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func createDatabaseDirectory() (string, error) {
	dbDir := os.Getenv("DB_DIR")
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return "", fmt.Errorf("could not create database dir: %w", err)
	}

	return filepath.Join(dbDir, "blkhell.db"), nil
}
