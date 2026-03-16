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
	"github.com/vague2k/blkhell/views/components"
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

func (h *Handler) HXReleaseProjectsTable(w http.ResponseWriter, r *http.Request) {
	projects, err := h.config.Database.GetProjectsByRelease(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	components.ProjectsTable(projects).Render(r.Context(), w)

	count := len(projects)
	if count <= 0 {
		fmt.Fprint(
			w,
			`<span id="release-projects-count" hx-swap-oob="true" class="text-muted-foreground text-xs">No projects to show yet</span>`,
		)
	} else {
		fmt.Fprintf(
			w,
			`<span id="release-projects-count" hx-swap-oob="true" class="text-muted-foreground text-xs">%d PROJECTS</span>`,
			count,
		)
	}
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

	_, err = h.config.Database.CreateRelease(r.Context(), database.CreateReleaseParams{
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

	toastSuccess(w, r, fmt.Sprintf("You made a new release for %s! Sick.", band.Name))
}

func (h *Handler) EditRelease(w http.ResponseWriter, r *http.Request) {
	releaseName := r.FormValue("release-name")
	releaseType := r.FormValue("release-type")
	releaseNum := r.FormValue("release-number")

	if releaseName == "" && releaseType == "" && releaseNum == "" {
		toastError(w, r, "At least 1 field has to be filled.")
		return
	}

	release, err := h.config.Database.GetReleaseByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to edit.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	params := database.UpdateReleaseParams{
		ID:     release.ID,
		Name:   release.Name,
		Type:   release.Type,
		Number: release.Number,
	}

	if releaseName != "" {
		params.Name = releaseName
	}
	if releaseType != "" {
		params.Type = releaseType
	}
	if releaseNum != "" {
		params.Number = releaseNum
	}

	_, err = h.config.Database.UpdateRelease(r.Context(), params)
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	toastSuccess(w, r, "Your changes have been saved! You may have to refresh the page to see your changes")
}

func (h *Handler) UploadReleaseAsset(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "Could not get user from context.")
		return
	}

	release, err := h.config.Database.GetReleaseByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to upload an asset for")
			return
		}
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	asset, err := h.FilesService.Upload(w, r, user.ID, release.ID, common.FileOwnerTypeRelease)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("'%s' was uploaded successfully!", asset.FullFilename()))
}

func (h *Handler) DeleteRelease(w http.ResponseWriter, r *http.Request) {
	release, err := h.config.Database.GetReleaseByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to upload an asset for")
			return
		}
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	err = h.config.Database.DeleteRelease(r.Context(), release.ID)
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	setToastCookie(w, "deleted_release_toast_message", fmt.Sprintf("The release %s was permanently deleted.", release.Name))
	w.Header().Set("HX-Redirect", "/band/"+release.BandID)
}
