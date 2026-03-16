package handlers

import (
	"github.com/vague2k/blkhell/views/pages"
	"net/http"
)

func (h *Handler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	toastCookieWarning(w, r, "deleted_band_toast_message")
	pages.Dashboard().Render(r.Context(), w)
}
