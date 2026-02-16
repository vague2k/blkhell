package handlers

import (
	"net/http"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.Auth.DestroySession(r)
	if err != nil {
		toastError(w, r, "500 Internal error: Could not destroy session.")
		return
	}

	w.Header().Set("HX-Redirect", "/login")
}
