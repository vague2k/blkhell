package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/handlers"
	"github.com/vague2k/blkhell/server/middleware"
)

func (s *Server) RegisterRoutes(h *handlers.Handler, middleware *middleware.Middleware) {
	s.router.With(middleware.RedirectIfAuth).Get("/login", h.LoginPage)
	s.router.Post("/login", h.Login)

	s.router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)

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

			// ---------- release pages ----------
			r.Route("/release", func(r chi.Router) {
				r.Get("/{id}", h.ReleasePage)
			})

			// ---------- project pages ----------
			r.Route("/project", func(r chi.Router) {
				r.Get("/{id}", h.ProjectPage)
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
				r.Put("/{id}/edit", h.EditBand)
				r.Delete("/{id}/delete", h.DeleteBand)
			})

			// ---------- release actions ----------
			r.Route("/release", func(r chi.Router) {
				r.Post("/create/{band-id}", h.CreateRelease)
				r.Put("/{id}/edit", h.EditRelease)
				r.Delete("/{id}/delete", h.DeleteRelease)
			})

			// ---------- release actions ----------
			r.Route("/project", func(r chi.Router) {
				r.Post("/create/{band-id}/{release-id}", h.CreateProject)
				r.Put("/{id}/edit", h.EditProject)
				r.Delete("/{id}/delete", h.DeleteProject)
			})

			// ---------- file actions ----------
			r.Route("/files", func(r chi.Router) {
				r.Get("/download/{id}", h.Download)
				r.Delete("/delete/{id}", h.Delete)
				r.Put("/edit/{id}", h.EditFile)

				// ---------- label file actions ----------
				r.Route("/label", func(r chi.Router) {
					r.Post("/upload", h.UploadLabelAsset)
				})

				// ---------- band file actions ----------
				r.Route("/bands", func(r chi.Router) {
					r.Post("/{id}/upload", h.UploadBandAsset)
				})

				// ---------- releases file actions ----------
				r.Route("/releases", func(r chi.Router) {
					r.Post("/{id}/upload", h.UploadReleaseAsset)
				})

				// ---------- projects file actions ----------
				r.Route("/projects", func(r chi.Router) {
					r.Post("/{id}/upload", h.UploadProjectAsset)
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
				r.Get("/chart", h.HXDashboardChart)
			})
			// ---------- label assets page fragments ----------
			r.Route("/label-assets", func(r chi.Router) {
				r.Get("/image-gallery", h.HXLabelAssetsImageGallery)
			})

			r.Route("/bands/{id}", func(r chi.Router) {
				r.Get("/image-gallery", h.HXBandsAssetsImageGallery)
				r.Get("/release-table", h.HXBandsReleaseTable)
				r.Get("/projects-table", h.HXBandProjectsTable)
				// r.Get("/{id}/search-image-gallery", )
			})

			r.Route("/releases/{id}", func(r chi.Router) {
				r.Get("/image-gallery", h.HXReleaseAssetsImageGallery)
				r.Get("/projects-table", h.HXReleaseProjectsTable)
			})

			r.Route("/projects/{id}", func(r chi.Router) {
				r.Get("/image-gallery", h.HXProjectsAssetsImageGallery)
			})
		})
	})
}
