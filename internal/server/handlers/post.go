package handlers

import (
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.Auth.Authenticate(r.Context(), username, password)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	sessionID, expires, err := h.Auth.CreateSession(r.Context(), user.ID)
	if err != nil {
		toastError(w, r, "500 Internal error: Could not create session.")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})

	w.Header().Set("HX-Redirect", "/dashboard")
}
