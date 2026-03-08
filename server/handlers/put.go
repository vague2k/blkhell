package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/database"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	oldPassword := r.FormValue("old-password")
	newPassword := r.FormValue("new-password")

	if username == "" && newPassword == "" {
		toastError(w, r, "You must enter a new username or password.")
		return
	}

	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "500 Internal error: Could not get current user.")
		return
	}

	// these will be left unchanged for now
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
			toastError(w, r, "500 Internal error: Could not hash new password.")
			return
		}
		params.PasswordHash = string(newPasswordHash)
	}

	_, err := h.DB.UpdateUser(r.Context(), params)
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
		if err := h.AuthService.DestroySession(w, r); err != nil {
			toastError(w, r, "500 Internal error: Could not destroy session.")
			return
		}

		toastSuccess(w, r, "Your changes have been saved! You must refresh the page and log in again with your new password.")
		return
	}

	// username only update toast
	toastSuccess(w, r, "You changes has been saved!")
}

func (h *Handler) EditBand(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("band-name")
	country := r.FormValue("band-country")

	// Make sure at least one field is set
	if name == "" && country == "" {
		toastError(w, r, "At least 1 field has to be filled.")
		return
	}

	band, err := h.DB.GetBandByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "This band doesn't exist.")
			return
		}

		toastError(w, r, "500 Internal error: Could not find band.")
		return
	}

	params := database.UpdateBandParams{
		ID:      band.ID,
		Name:    band.Name,
		Country: band.Country,
	}

	if name != "" {
		params.Name = name
	}
	if country != "" {
		params.Country = country
	}

	_, err = h.DB.UpdateBand(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "The band you're trying to update doesn't exist.")
			return
		}

		toastError(w, r, "500 Internal error: Could not update band.")
		return
	}

	toastSuccess(w, r, "Your changes have been saved!")
}
