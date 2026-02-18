package handlers

import (
	"net/http"

	"github.com/vague2k/blkhell/views/components"
	"github.com/vague2k/blkhell/views/pages"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.GetUserFromRequest(r)
	if err == nil && user != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	pages.Login().Render(r.Context(), w)
}

func (h *Handler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, "Couldn't get user from request", http.StatusInternalServerError)
		return
	}
	pages.Dashboard(user).Render(r.Context(), w)
}

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, "Couldn't get user from request", http.StatusInternalServerError)
		return
	}
	pages.Settings(user).Render(r.Context(), w)
}

func (h *Handler) HXUploadedImages(w http.ResponseWriter, r *http.Request) {
	images, err := h.DB.GetImages(r.Context())
	if err != nil {
		http.Error(w, "Couldn't get user from request", http.StatusInternalServerError)
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)
}
