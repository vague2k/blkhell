package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/vague2k/blkhell/server/database"
	_ "modernc.org/sqlite"
)

const (
	EnvDevelopment = "development"
	EnvTesting     = "testing"
	EnvProduction  = "production"
)

type Config struct {
	Environment string
	Port        string
	UploadsDir  string
	DBDir       string

	SqlDB    *sql.DB
	Database *database.Queries
}

func Init() *Config {
	if err := requireFFmpeg(); err != nil {
		log.Fatal("Fatal config error: missing ffprobe (install ffmpeg)")
	}

	if err := loadEnvVars(); err != nil {
		log.Fatalf("Fatal config error: loading env vars: %v", err)
	}

	c := &Config{
		Environment: runtimeEnv(),
		Port:        requireEnv("PORT"),
		UploadsDir:  requireEnv("UPLOADS_DIR"),
		DBDir:       requireEnv("DB_DIR"),
	}

	if err := os.MkdirAll(c.UploadsDir, 0o755); err != nil {
		log.Fatalf("Fatal config error: creating uploads dir: %v", err)
	}

	db, err := c.openDB()
	if err != nil {
		log.Fatalf("Fatal config error: opening database: %v", err)
	}

	c.SqlDB = db
	c.Database = database.New(db)

	return c
}

func (c *Config) openDB() (*sql.DB, error) {
	connStr, err := c.dbConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	return db, nil
}

func (c *Config) dbConnectionString() (string, error) {
	switch c.Environment {

	case EnvTesting:
		return ":memory:", nil

	case EnvDevelopment, EnvProduction:
		if err := os.MkdirAll(c.DBDir, 0o755); err != nil {
			return "", fmt.Errorf("creating database dir: %w", err)
		}

		return filepath.Join(c.DBDir, "blkhell.db"), nil

	default:
		return "", fmt.Errorf("unknown environment: %s", c.Environment)
	}
}

func loadEnvVars() error {
	if runtimeEnv() == EnvTesting {
		return nil
	}
	return godotenv.Load()
}

func runtimeEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		return "unknown"
	}
	return env
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}

func requireFFmpeg() error {
	// ffprobe comes with ffmpeg
	_, err := exec.LookPath("ffprobe")
	return err
}
