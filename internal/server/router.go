package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/internal/server/auth"
	"github.com/vague2k/blkhell/internal/server/handlers"
)

func (s *Server) RegisterRoutes(authService *auth.Service) {
	h := handlers.NewHandler(authService)

	// pages
	s.router.Get("/login", h.LoginPage)

	// all pages that require auth to access
	s.router.Group(func(r chi.Router) {
		r.Use(authService.RequireAuth)
		r.Get("/", h.DashboardPage)
		r.Get("/dashboard", h.DashboardPage)
	})

	// backend endpoints
	s.router.Post("/login", h.Login)
}
