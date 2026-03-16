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

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	projectName := r.FormValue("project-name")
	projectType := r.FormValue("project-type")

	if projectName == "" {
		toastError(w, r, "'Project Name' is required")
		return
	}
	if projectType == "" {
		toastError(w, r, "'Project Type' is required")
		return
	}

	band, err := h.config.Database.GetBandByID(r.Context(), chi.URLParam(r, "band-id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band to create project for.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	release, err := h.config.Database.GetReleaseByID(r.Context(), chi.URLParam(r, "release-id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to create project for.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	project, err := h.config.Database.CreateProject(r.Context(), database.CreateProjectParams{
		ID:        uuid.NewString(),
		BandID:    band.ID,
		ReleaseID: release.ID,
		Name:      projectName,
		Type:      projectType,
		Status:    common.ProjectStatusPending,
	})
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("You new project '%s' has been added created for %s's release '%s'! beast.", project.Name, band.Name, release.Name))
}

func (h *Handler) EditProject(w http.ResponseWriter, r *http.Request) {
	projectName := r.FormValue("project-name")
	projectType := r.FormValue("project-type")
	projectStatus := r.FormValue("project-status")

	if projectName == "" && projectType == "" && projectStatus == "" {
		toastError(w, r, "At least 1 field has to be filled.")
		return
	}

	project, err := h.config.Database.GetProjectByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to edit.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	params := database.UpdateProjectParams{
		ID:     project.ID,
		Name:   project.Name,
		Type:   project.Type,
		Status: project.Status,
	}

	if projectName != "" {
		params.Name = projectName
	}
	if projectType != "" {
		params.Type = projectType
	}
	if projectStatus != "" {
		params.Status = projectStatus
	}

	_, err = h.config.Database.UpdateProject(r.Context(), params)
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	toastSuccess(w, r, "Your changes have been saved! You may have to refresh the page to see your changes")
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

	asset, err := h.FilesService.Upload(w, r, user.ID, project.ID, common.FileOwnerTypeProject)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("'%s' was uploaded successfully!", asset.FullFilename()))
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	project, err := h.config.Database.GetProjectByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to upload an asset for")
			return
		}
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	err = h.config.Database.DeleteProject(r.Context(), project.ID)
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	setToastCookie(w, "deleted_project_toast_message", fmt.Sprintf("The project %s was permanently deleted.", project.Name))
	w.Header().Set("HX-Redirect", "/band/"+project.BandID)
}
