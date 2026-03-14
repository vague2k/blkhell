package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/data"
	serverErrors "github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/pages"
)

func (h *Handler) ProjectPage(w http.ResponseWriter, r *http.Request) {
	project, err := h.config.Database.GetProjectByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to display.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	release, err := h.config.Database.GetReleaseByID(r.Context(), project.ReleaseID)
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

	pages.Project(&band, &release, &project).Render(r.Context(), w)
}

func (h *Handler) UploadProjectAsset(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "Could not get user from context.")
		return
	}

	project, err := h.config.Database.GetProjectByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get project to upload an asset for")
			return
		}
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	asset, err := h.FilesService.Upload(w, r, user.ID, project.ID, data.FileOwnerTypeProject)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("'%s' was uploaded successfully!", asset.FullFilename()))
}
