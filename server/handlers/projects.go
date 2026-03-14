package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	serverErrors "github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/components"
)

func (h *Handler) HXReleaseProjectsTable(w http.ResponseWriter, r *http.Request) {
	projects, err := h.config.Database.GetProjectsByRelease(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	components.ProjectsTable(projects).Render(r.Context(), w)
	if len(projects) > 0 {
		fmt.Fprintf(
			w,
			`<span id="release-projects-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d PROJECTS</span>`,
			len(projects),
		)
	}
}
