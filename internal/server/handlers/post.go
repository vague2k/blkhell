package handlers

import (
	"net/http"

	"github.com/vague2k/blkhell/views/templui/toast"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.auth.Authenticate(r.Context(), username, password)
	if err != nil {
		toast.Toast(toast.Props{
			Icon:          true,
			Title:         "Invalid credentials",
			Description:   err.Error(),
			Variant:       toast.VariantError,
			Position:      toast.PositionTopRight,
			Dismissible:   true,
			ShowIndicator: true,
		}).Render(r.Context(), w)
		return
	}

	sessionID, expires, err := h.auth.CreateSession(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
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
