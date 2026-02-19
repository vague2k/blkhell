package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vague2k/blkhell/server"
	"github.com/vague2k/blkhell/server/auth"
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/server/handlers"
)

func main() {
	err := godotenv.Load(".env", "version")
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	queries, err := database.Init()
	if err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	authService := auth.New(queries)

	handler := handlers.NewHandler(authService, queries)

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT env var is not set")
	}

	s := server.NewServer(port)

	s.SetupAssetsRoutes()
	s.SetupUploadRoutes()
	s.RegisterRoutes(handler)

	if err := s.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
