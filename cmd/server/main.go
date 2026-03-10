package main

import (
	"log"

	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/server"
	"github.com/vague2k/blkhell/server/handlers"
	"github.com/vague2k/blkhell/server/middleware"
)

func main() {
	cfg := config.Init()
	handler := handlers.NewHandler(cfg.Database)
	middleware := middleware.New(cfg.Database)

	s := server.NewServer(cfg.Port)
	s.SetupAssetsRoutes()
	s.SetupUploadRoutes()
	s.RegisterRoutes(handler, middleware)

	if err := s.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
