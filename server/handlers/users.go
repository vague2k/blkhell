package handlers

import (
	"net/http"

	"github.com/vague2k/blkhell/server/database"
	serverErrors "github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/pages"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	pages.Settings().Render(r.Context(), w)
}

func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	oldPassword := r.FormValue("old-password")
	newPassword := r.FormValue("new-password")

	if username == "" && newPassword == "" {
		toastError(w, r, "You must enter a new username or password.")
		return
	}

	user, ok := h.authService.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "Could not get user from context.")
		return
	}

	params := database.UpdateUserParams{
		ID:           user.ID,
		Role:         user.Role,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}

	if username != "" {
		if username == user.Username {
			toastError(w, r, "New username must be different from the current username.")
			return
		}
		params.Username = username
	}

	if newPassword != "" {
		if oldPassword == "" {
			toastError(w, r, "You must enter your old password to set a new one.")
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)) != nil {
			toastError(w, r, "The old password you entered is incorrect, try again")
			return
		}

		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			toastError(w, r, "Could not create new password.")
			return
		}
		params.PasswordHash = string(newPasswordHash)
	}

	_, err := h.config.Database.UpdateUser(r.Context(), params)
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	if newPassword != "" {
		if err := h.authService.DestroySession(w, r); err != nil {
			toastError(w, r, err.Error())
			return
		}

		setToastCookie(w, "changed_password_toast_message", "You set a new password! You must log in again.")
		w.Header().Set("HX-Redirect", "/login")
		return
	}

	toastSuccess(w, r, "You changes has been saved!")
}
