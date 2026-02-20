package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/handlers"
)

func (s *Server) RegisterRoutes(h *handlers.Handler) {
	s.router.With(h.Auth.RedirectIfAuth).Get("/login", h.LoginPage)
	s.router.Post("/login", h.Login)

	s.router.Group(func(r chi.Router) {
		r.Use(h.Auth.RequireAuth)
		r.Get("/", h.DashboardPage)
		r.Get("/dashboard", h.DashboardPage)
		r.Get("/settings", h.SettingsPage)
		r.Post("/upload", h.Upload)
		r.Delete("/logout", h.Logout)
		r.Delete("/images/{id}", h.DeleteImage)
		r.Put("/edit", h.Edit)
		r.Get("/hx/uploaded-images", h.HXUploadedImages)
		r.Get("/hx/search-image", h.HXSearchImages)
		r.Get("/images/download/{id}", h.HXDownloadImage)
	})
}
