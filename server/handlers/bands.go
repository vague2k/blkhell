package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/vague2k/blkhell/common"
	"github.com/vague2k/blkhell/server/database"
	serverErrors "github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/pages"
)

func (h *Handler) BandPage(w http.ResponseWriter, r *http.Request) {
	band, err := h.config.Database.GetBandByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band to display.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	toastCookieWarning(w, r, "deleted_release_toast_message")
	toastCookieWarning(w, r, "deleted_project_toast_message")
	pages.Band(&band).Render(r.Context(), w)
}

func (h *Handler) CreateBand(w http.ResponseWriter, r *http.Request) {
	bandName := r.FormValue("band-name")
	bandCountry := r.FormValue("band-country")

	if bandName == "" {
		toastError(w, r, "'Band name' is required")
		return
	}
	if bandCountry == "" {
		toastError(w, r, "'Band country' is required")
		return
	}

	band, err := h.config.Database.CreateBand(r.Context(), database.CreateBandParams{
		ID:      uuid.NewString(),
		Name:    bandName,
		Country: bandCountry,
	})
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("You new band '%s' has been added to the roster! beast.", band.Name))
}

func (h *Handler) EditBand(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("band-name")
	country := r.FormValue("band-country")

	if name == "" && country == "" {
		toastError(w, r, "At least 1 field has to be filled.")
		return
	}

	band, err := h.config.Database.GetBandByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band to edit.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
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

	_, err = h.config.Database.UpdateBand(r.Context(), params)
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	toastSuccess(w, r, "Your changes have been saved! You may have to refresh the page to see your changes")
}

func (h *Handler) UploadBandAsset(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "Could not get user from context.")
		return
	}

	band, err := h.config.Database.GetBandByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band to upload an asset for")
			return
		}
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	asset, err := h.FilesService.Upload(w, r, user.ID, band.ID, common.FileOwnerTypeBand)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("'%s' was uploaded successfully!", asset.FullFilename()))
}

func (h *Handler) DeleteBand(w http.ResponseWriter, r *http.Request) {
	band, err := h.config.Database.GetBandByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band to delete")
			return
		}
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	err = h.config.Database.DeleteBand(r.Context(), band.ID)
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	setToastCookie(w, "deleted_band_toast_message", fmt.Sprintf("The band %s was permanently deleted.", band.Name))
	w.Header().Set("HX-Redirect", "/")
}
