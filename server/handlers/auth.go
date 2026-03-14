package handlers

import (
	"net/http"

	"github.com/vague2k/blkhell/views/pages"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if ok && user != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	toastCookieSuccess(w, r, "changed_password_toast_message")
	pages.Login().Render(r.Context(), w)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.AuthService.Authenticate(r.Context(), username, password)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	sessionToken, expires, err := h.AuthService.CreateSession(r.Context(), user.ID)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})

	w.Header().Set("HX-Redirect", "/dashboard")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.AuthService.DestroySession(w, r)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	w.Header().Set("HX-Redirect", "/login")
}
