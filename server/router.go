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
			r.Get("/label-assets", h.LabelAssetsPage)
			r.Get("/settings", h.SettingsPage)

			// ---------- band pages ----------
			r.Route("/band", func(r chi.Router) {
				r.Get("/{id}", h.BandPage)
			})
		})

		// ---------- actions ----------
		r.Route("/actions", func(r chi.Router) {
			// ---------- user actions ----------
			r.Route("/users", func(r chi.Router) {
				r.Delete("/logout", h.Logout)
				r.Put("/edit", h.EditUser)
			})

			// ---------- band actions ----------
			r.Route("/bands", func(r chi.Router) {
				r.Post("/create", h.CreateBand)
			})

			// ---------- file actions ----------
			r.Route("/files", func(r chi.Router) {
				r.Get("/download/{id}", h.Download)
				r.Delete("/delete/{id}", h.Delete)

				// ---------- label file actions ----------
				r.Route("/label", func(r chi.Router) {
					r.Post("/upload", h.UploadLabelAsset)
					r.Delete("/delete/{id}", h.Delete)
				})

				// ---------- band file actions ----------
				r.Route("/bands", func(r chi.Router) {
					r.Post("/{id}/upload", h.UploadBandAsset)
				})
			})
		})

		// ---------- htmx fragments ----------
		r.Route("/hx", func(r chi.Router) {
			// ---------- dashboard page fragments ----------
			r.Route("/sidebar", func(r chi.Router) {
				r.Get("/user-dropdown", h.HXSidebarUserDropdown)
				r.Get("/bands-dropdown", h.HXSidebarBandsDropdown)
			})
			// ---------- dashboard page fragments ----------
			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/cards", h.HXDashboardCards)
				r.Get("/table", h.HXDashboardTable)
			})
			// ---------- label assets page fragments ----------
			r.Route("/label-assets", func(r chi.Router) {
				r.Get("/image-gallery", h.HXImageGallery)
				r.Get("/search-image-gallery", h.HXSearchImageGallery)
			})
		})
	})
}
