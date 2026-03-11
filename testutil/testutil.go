package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vague2k/blkhell/config"
	_ "modernc.org/sqlite"
)

type TB interface {
	Fatalf(format string, args ...interface{})
	Cleanup(func())
}

func NewTestConfig(t TB) (*config.Config, func()) {
	os.Setenv("GO_ENV", "testing")
	os.Setenv("PORT", "8080")
	os.Setenv("UPLOADS_DIR", "/tmp/blkhell_test_uploads")
	os.Setenv("DB_DIR", "/tmp/blkhell_test_database")

	cfg := config.Init()

	if err := migrateUpTest(cfg.SqlDB); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		cfg.SqlDB.Close()
		os.RemoveAll("/tmp/blkhell_test_uploads")
		os.RemoveAll("/tmp/blkhell_test_database")
	}
	return cfg, cleanup
}

func migrateUpTest(db *sql.DB) error {
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))
	migrationsPath := filepath.Join(projectRoot, "server", "database", "migrations")

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite",
		driver,
	)
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

func Context() context.Context {
	return context.Background()
}

func RenderComponent(component templ.Component) (*goquery.Document, error) {
	reader, writer := io.Pipe()
	go func() {
		err := component.Render(context.Background(), writer)
		if err != nil {
			writer.CloseWithError(err)
			return
		}
		writer.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
