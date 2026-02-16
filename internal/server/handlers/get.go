package handlers

import (
	"net/http"

	"github.com/vague2k/blkhell/views/pages"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromRequest(r)
	if err == nil && user != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	pages.Login().Render(r.Context(), w)
}

func (h *Handler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	pages.Dashboard().Render(r.Context(), w)
}
