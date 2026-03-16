package testutil

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/server/handlers"
	"github.com/vague2k/blkhell/server/services"
	_ "modernc.org/sqlite"
)

type TB interface {
	Fatalf(format string, args ...interface{})
	Cleanup(func())
}

type Test struct {
	AuthService      *services.AuthService
	FilesService     *services.FilesService
	DashboardService *services.DashboardService
	Config           *config.Config
	Handler          *handlers.Handler

	ctx context.Context
	t   TB // *testing.T
}

func NewTest(t TB) *Test {
	os.Setenv("GO_ENV", "testing")
	os.Setenv("PORT", "8080")
	os.Setenv("UPLOADS_DIR", "/tmp/blkhell_test_uploads")
	os.Setenv("DB_DIR", "/tmp/blkhell_test_database")

	cfg := config.Init()

	if err := migrateUpTest(cfg.SqlDB); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		if err := cfg.SqlDB.Close(); err != nil {
			t.Fatalf("failed to close db: %v", err)
		}
		os.RemoveAll("/tmp/blkhell_test_uploads")
		os.RemoveAll("/tmp/blkhell_test_database")
	}

	t.Cleanup(cleanup)

	return &Test{
		AuthService:      services.NewAuthService(cfg),
		FilesService:     services.NewFilesService(cfg),
		DashboardService: services.NewDashboardService(cfg),
		Handler:          handlers.NewHandler(cfg),
		Config:           cfg,
		ctx:              context.Background(),
		t:                t,
	}
}

func (t *Test) Context() context.Context {
	return t.ctx
}

func (t *Test) NewRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func (t *Test) NewFormRequest(path string, values url.Values) *http.Request {
	return httptest.NewRequestWithContext(t.ctx, http.MethodPost, path, strings.NewReader(values.Encode()))
}

func (t *Test) RenderComponent(component templ.Component) (*goquery.Document, error) {
	reader, writer := io.Pipe()
	go func() {
		err := component.Render(t.ctx, writer)
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

func (t *Test) RandomUsername() string {
	return "user_" + randomString(8)
}

func (t *Test) RandomPassword() string {
	return "pass_" + randomString(12)
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

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}
