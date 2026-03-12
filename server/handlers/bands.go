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
	"github.com/vague2k/blkhell/views/components"
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
	pages.Band(&band).Render(r.Context(), w)
}

func (h *Handler) HXBandsReleaseTable(w http.ResponseWriter, r *http.Request) {
	releases, err := h.config.Database.GetReleasesByBand(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	components.BandReleasesTable(releases).Render(r.Context(), w)
	fmt.Fprintf(
		w,
		`<span id="band-releases-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d RELEASES</span>`,
		len(releases),
	)
}

func (h *Handler) HXBandProjectsTable(w http.ResponseWriter, r *http.Request) {
	projects, err := h.config.Database.GetProjectsByBandID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	components.BandProjectsTable(projects).Render(r.Context(), w)
	if len(projects) > 0 {
		fmt.Fprintf(
			w,
			`<span id="band-projects-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d PROJECTS</span>`,
			len(projects),
		)
	}
}

func (h *Handler) CreateBand(w http.ResponseWriter, r *http.Request) {
	bandName := r.FormValue("band-name")
	bandCountry := r.FormValue("band-country")
	releaseName := r.FormValue("release-name")
	releaseType := r.FormValue("release-type")
	releaseNum := r.FormValue("release-number")
	projectName := r.FormValue("project-name")
	projectType := r.FormValue("project-type")

	switch true {
	case bandName == "":
		toastError(w, r, "'Band name' is required")
		return
	case releaseName == "":
		toastError(w, r, "'Release Name' is required")
		return
	case releaseType == "":
		toastError(w, r, "'Release Type' is required")
		return
	case releaseNum == "":
		toastError(w, r, "'Release No.' is required")
		return
	case bandCountry == "":
		toastError(w, r, "'Band country' is required")
		return
	}

	if (projectName == "" && projectType != "") || (projectName != "" && projectType == "") {
		toastError(w, r, "'Project Name' and 'Project Type' must both be filled or both left empty")
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

	toastSuccess(w, r, "Your changes have been saved!")
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

	asset, err := h.FilesService.Upload(r, user.ID, band.ID, "band")
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("'%s' was uploaded successfully!", asset.FullFilename()))
}
