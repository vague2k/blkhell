package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

// MigrateUp applies all pending up migrations
func MigrateUp(db *sql.DB) error {
	m, err := migrator(db)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	fmt.Printf("Current migration version: %d (dirty: %v)\n", version, dirty)

	fmt.Println("Running migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	newVersion, _, _ := m.Version()
	if newVersion == version {
		fmt.Println("No migrations to apply")
	} else {
		fmt.Printf("Successfully applied migrations (now at version %d)\n", newVersion)
	}
	return nil
}

// MigrateDown applies one down migration
func MigrateDown(db *sql.DB) error {
	m, err := migrator(db)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	fmt.Printf("Current migration version: %d (dirty: %v)\n", version, dirty)

	fmt.Println("Rolling back one migration...")
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rollback failed: %w", err)
	}

	newVersion, _, _ := m.Version()
	fmt.Printf("Successfully rolled back migration (now at version %d)\n", newVersion)
	return nil
}

func MigrateForce(db *sql.DB, version int) error {
	m, err := migrator(db)
	if err != nil {
		return err
	}
	return m.Force(version)
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
