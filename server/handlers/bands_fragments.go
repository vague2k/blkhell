package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	serverErrors "github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/components"
)

func (h *Handler) HXBandsAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.config.Database.GetBandImageFilesByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "Could not get band images to display")
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)

	count := len(images)
	if count <= 0 {
		fmt.Fprint(
			w,
			`<span id="band-assets-image-count" hx-swap-oob="true"></span>`,
		)
	} else {
		fmt.Fprintf(
			w,
			`<span id="band-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
			count,
		)
	}
}

func (h *Handler) HXBandsReleaseTable(w http.ResponseWriter, r *http.Request) {
	releases, err := h.config.Database.GetReleasesByBand(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	components.ReleasesTable(releases).Render(r.Context(), w)

	count := len(releases)
	if count <= 0 {
		fmt.Fprint(
			w,
			`<span id="band-releases-count" hx-swap-oob="true" class="text-muted-foreground text-xs">No releases to show yet</span>`,
		)
	} else {
		fmt.Fprintf(
			w,
			`<span id="band-releases-count" hx-swap-oob="true" class="text-muted-foreground text-xs">%d RELEASES</span>`,
			count,
		)
	}
}

func (h *Handler) HXBandProjectsTable(w http.ResponseWriter, r *http.Request) {
	projects, err := h.config.Database.GetProjectsByBandID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	components.ProjectsTable(projects).Render(r.Context(), w)
	count := len(projects)
	if count <= 0 {
		fmt.Fprint(
			w,
			`<span id="band-projects-count" hx-swap-oob="true" class="text-muted-foreground text-xs">No projects to show yet</span>`,
		)
	} else {
		fmt.Fprintf(
			w,
			`<span id="band-projects-count" hx-swap-oob="true" class="text-muted-foreground text-xs">%d PROJECTS</span>`,
			count,
		)
	}
}
