package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vague2k/blkhell/server/database"
)

type Config struct {
	Environment string
	Port        string
	Database    *database.Queries
	UploadsDir  string
}

func Init() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env: %v\n", err)
	}

	queries, err := database.Init()
	if err != nil {
		log.Fatalf("Database init failed: %v\n", err)
	}

	return &Config{
		Environment: getEnv("GO_ENV"),
		Port:        getEnv("PORT"),
		Database:    queries,
		UploadsDir:  getEnv("UPLOADS_DIR"),
	}
}

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Error: '%s' env var is not set\n", key)
	}
	return v
}
