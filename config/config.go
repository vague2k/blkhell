package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

	SqlDB    *sql.DB
	Database *database.Queries
}

func Init() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env: %v\n", err)
	}

	c := &Config{
		Environment: getEnv("GO_ENV"),
		Port:        getEnv("PORT"),
		UploadsDir:  getEnv("UPLOADS_DIR"),
	}

	err = c.createUploadRootDir()
	if err != nil {
		log.Fatalf("Error: failed to create upload root directory: %v\n", err)
	}

	db, err := c.openDB()
	if err != nil {
		log.Fatalf("Error: failed to open database: %v\n", err)
	}

	c.SqlDB = db
	c.Database = database.New(db)

	return c
}

func (c *Config) openDB() (*sql.DB, error) {
	connStr, err := c.dbConnectionString()
	if err != nil {
		return nil, fmt.Errorf("could not create db connection string: %w", err)
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
		dbDir := getEnv("DB_DIR")
		if err := os.MkdirAll(dbDir, 0o755); err != nil {
			return "", fmt.Errorf("could not create database dir: %w", err)
		}
		return filepath.Join(dbDir, "blkhell.db"), nil

	default:
		return "", fmt.Errorf("unknown environment: %s", c.Environment)
	}
}

func (c *Config) createUploadRootDir() error {
	if err := os.MkdirAll(c.UploadsDir, 0o755); err != nil {
		return err
	}
	return nil
}

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Error: '%s' env var is not set\n", key)
	}
	return v
}
