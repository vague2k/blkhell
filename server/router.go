package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/auth"
	"github.com/vague2k/blkhell/server/handlers"
)

func (s *Server) RegisterRoutes(authService *auth.Service) {
	h := handlers.NewHandler(authService)

	// pages
	s.router.Handle("/login", h.LoginPage())
	// all pages that require auth to access
	s.router.Group(func(r chi.Router) {
		r.Use(authService.RequireAuth)
		r.Handle("/", h.DashboardPage())
		r.Handle("/dashboard", h.DashboardPage())
	})

	// endpoints
	s.router.Post("/login", h.Login)
}
