package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/handlers"
)

func (s *Server) RegisterRoutes(h *handlers.Handler) {
	s.router.Group(func(r chi.Router) {
		r.Use(h.Auth.RequireAuth)
		r.Get("/login", h.LoginPage)
		r.Get("/", h.DashboardPage)
		r.Get("/dashboard", h.DashboardPage)
		r.Get("/settings", h.SettingsPage)
	})

	// backend endpoints
	s.router.Post("/login", h.Login)
	s.router.Post("/upload", h.Upload)
	s.router.Delete("/logout", h.Logout)
	s.router.Delete("/images/{id}", h.DeleteImage)
	s.router.Put("/edit", h.Edit)
	s.router.Get("/hx/uploaded-images", h.HXUploadedImages)
	s.router.Get("/hx/search-image", h.HXSearchImages)
	s.router.Get("/images/download/{id}", h.HXDownloadImage)
}
