package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/vague2k/blkhell/server/data"
	"github.com/vague2k/blkhell/server/database"
	serverErrors "github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/pages"
)

func (h *Handler) ReleasePage(w http.ResponseWriter, r *http.Request) {
	release, err := h.config.Database.GetReleaseByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to display.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}
	band, err := h.config.Database.GetBandByID(r.Context(), release.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band from release.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	pages.Release(&band, &release).Render(r.Context(), w)
}

func (h *Handler) CreateRelease(w http.ResponseWriter, r *http.Request) {
	band, err := h.config.Database.GetBandByID(r.Context(), chi.URLParam(r, "band-id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band to create a release for")
			return
		}
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	releaseName := r.FormValue("release-name")
	releaseType := r.FormValue("release-type")
	releaseNum := r.FormValue("release-number")
	releaseSongCount := r.FormValue("release-song-count")

	switch true {
	case releaseName == "":
		toastError(w, r, "'Release Name' is required")
		return
	case releaseType == "":
		toastError(w, r, "'Release Type' is required")
		return
	case releaseNum == "":
		toastError(w, r, "'Release No.' is required")
		return
	case releaseSongCount == "":
		toastError(w, r, "'Amount of Songs' is required")
		return
	}

	songCount, err := strconv.ParseInt(releaseSongCount, 0, 64)
	if err != nil {
		toastError(w, r, "'Amount of Songs' input is not a valid number")
		return
	}

	if songCount <= 0 {
		toastError(w, r, "A release can't have 0 or negative songs")
		return
	}

	if releaseType == data.ReleaseTypeSingle && songCount > 1 {
		toastError(w, r, "A single can only have 1 song")
		return
	}

	_, err = h.config.Database.CreateRelease(r.Context(), database.CreateReleaseParams{
		ID:     uuid.NewString(),
		BandID: band.ID,
		Name:   releaseName,
		Type:   releaseType,
		Number: releaseNum,
		// SongCount: songCount,
	})
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("Your new band '%s' has been added to the roster! beast.", band.Name))
}
