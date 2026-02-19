package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/vague2k/blkhell/server/database"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	oldPassword := r.FormValue("old-password")
	newPassword := r.FormValue("new-password")

	if username == "" && newPassword == "" {
		toastError(w, r, "You must enter a new username or password.")
		return
	}

	user, err := h.Auth.GetUserFromRequest(r)
	if err != nil {
		toastError(w, r, "500 Internal error: Could not get current user.")
		return
	}

	// these will be left unchanged for now
	params := database.UpdateUserParams{
		ID:   user.ID,
		Role: user.Role,
	}

	if username != "" {
		if username == user.Username {
			toastError(w, r, "New username must be different from the current username.")
			return
		}
		params.Username = username
	} else {
		params.Username = user.Username
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
			toastError(w, r, "500 Internal error: Could not hash new password.")
			return
		}
		params.PasswordHash = string(newPasswordHash)
	} else {
		params.PasswordHash = user.PasswordHash
	}

	_, err = h.DB.UpdateUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "The user you're trying to update doesn't exist.")
			return
		}

		toastError(w, r, "500 Internal error: Could not update user.")
		return
	}

	// destroy session if password changed
	if newPassword != "" {
		if err := h.Auth.DestroySession(r); err != nil {
			toastError(w, r, "500 Internal error: Could not destroy session.")
			return
		}

		toastSuccess(w, r, "Your changes have been saved! You must refresh the page and log in again with your new password.")
		return
	}

	// username only update toast
	toastSuccess(w, r, "You changes has been saved!")
}
