package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/handlers"
)

func (s *Server) RegisterRoutes(h *handlers.Handler) {
	s.router.With(h.Middleware.RedirectIfAuth).Get("/login", h.LoginPage)
	s.router.Post("/login", h.Login)

	s.router.Group(func(r chi.Router) {
		r.Use(h.Middleware.RequireAuth)

		// ---------- pages ----------
		r.Route("/", func(r chi.Router) {
			r.Get("/", h.DashboardPage)
			r.Get("/dashboard", h.DashboardPage)
			r.Get("/settings", h.SettingsPage)
		})

		// ---------- actions ----------
		r.Route("/actions", func(r chi.Router) {
			// ---------- user actions ----------
			r.Route("/users", func(r chi.Router) {
				r.Delete("/logout", h.Logout)
				r.Put("/edit", h.EditUser)
			})

			// ---------- image actions ----------
			r.Route("/images", func(r chi.Router) {
				r.Post("/upload", h.Upload)
				r.Delete("/delete/{id}", h.DeleteImage)
				r.Get("/download/{id}", h.Download)
			})
		})

		// ---------- htmx fragments ----------
		r.Route("/hx", func(r chi.Router) {
			// ---------- dashboard page fragments ----------
			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/image-gallery", h.HXImageGallery)
				r.Get("/search-image-gallery", h.HXSearchImageGallery)
			})
		})
	})
}
