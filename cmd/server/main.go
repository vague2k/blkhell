package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vague2k/blkhell/server"
	"github.com/vague2k/blkhell/server/auth"
	"github.com/vague2k/blkhell/server/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	queries, err := database.Init()
	if err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	authService := auth.New(queries)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	s := server.NewServer(port)

	s.SetupAssetsRoutes()
	s.RegisterRoutes(authService)

	if err := s.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
