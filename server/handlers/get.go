package handlers

import (
	"fmt"
	"net/http"

	"github.com/vague2k/blkhell/server/database"
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
	// render image count (hx-oob-swap)
	fmt.Fprintf(
		w,
		`<span id="dashboard-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d IMAGES</span>`,
		len(images),
	)
}

func (h *Handler) HXSearchImages(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("q")
	images, err := h.DB.GetImagesByPartialName(r.Context(), database.GetImagesByPartialNameParams{
		Filename: "%" + input + "%",
		Ext:      "%" + input + "%",
	})
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	components.ImageGallery(images).Render(r.Context(), w)
	// render image count (hx-oob-swap)
	fmt.Fprintf(
		w,
		`<span id="dashboard-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d IMAGES</span>`,
		len(images),
	)
}
