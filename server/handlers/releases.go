package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	projectName := r.FormValue("project-name")
	projectType := r.FormValue("project-type")

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
	}

	if (projectName == "" && projectType != "") || (projectName != "" && projectType == "") {
		toastError(w, r, "'Project Name' and 'Project Type' must both be filled or both left empty")
		return
	}

	release, err := h.config.Database.CreateRelease(r.Context(), database.CreateReleaseParams{
		ID:     uuid.NewString(),
		BandID: band.ID,
		Name:   releaseName,
		Type:   releaseType,
		Number: releaseNum,
	})
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	if projectName != "" && projectType != "" {
		_, err := h.config.Database.CreateProject(r.Context(), database.CreateProjectParams{
			ID:        uuid.NewString(),
			BandID:    band.ID,
			ReleaseID: release.ID,
			Name:      projectName,
			Type:      projectType,
		})
		if err != nil {
			toastError(w, r, serverErrors.ErrDb.Error())
			return
		}
	}

	toastSuccess(w, r, fmt.Sprintf("Your new band '%s' has been added to the roster! beast.", band.Name))
}
