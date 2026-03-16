package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/views/components"
)

func (h *Handler) HXProjectsAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.config.Database.GetProjectImageFilesByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "Could not get releases images to display")
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)

	count := len(images)
	if count <= 0 {
		fmt.Fprint(
			w,
			`<span id="project-assets-image-count" hx-swap-oob="true"></span>`,
		)
	} else {
		fmt.Fprintf(
			w,
			`<span id="project-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
			count,
		)
	}
}
