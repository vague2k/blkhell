package main

import (
	"log"

	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/server"
	"github.com/vague2k/blkhell/server/handlers"
)

func main() {
	cfg := config.Init()
	handler := handlers.NewHandler(cfg.Database)

	s := server.NewServer(cfg.Port)
	s.SetupAssetsRoutes()
	s.SetupUploadRoutes()
	s.RegisterRoutes(handler)

	if err := s.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
