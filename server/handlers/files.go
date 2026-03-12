package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/views/components"
	"github.com/vague2k/blkhell/views/pages"
)

func (h *Handler) LabelAssetsPage(w http.ResponseWriter, r *http.Request) {
	pages.LabelAssets().Render(r.Context(), w)
}

func (h *Handler) HXLabelAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.config.Database.GetLabelImageFiles(r.Context())
	if err != nil {
		toastError(w, r, "Could not get label files.")
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)
	fmt.Fprintf(
		w,
		`<span id="label-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
		len(images),
	)
}

func (h *Handler) HXBandsAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.config.Database.GetBandImageFilesByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "Could not get band images to display")
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)
	fmt.Fprintf(
		w,
		`<span id="band-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
		len(images),
	)
}

func (h *Handler) HXSearchLabelAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("q")
	images, err := h.config.Database.GetFileByPartialName(r.Context(), database.GetFileByPartialNameParams{
		Filename: "%" + input + "%",
		Ext:      "%" + input + "%",
	})
	if err != nil {
		toastError(w, r, "Search for asset image failed")
		return
	}

	components.ImageGallery(images).Render(r.Context(), w)
	fmt.Fprintf(
		w,
		`<span id="label-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
		len(images),
	)
}

func (h *Handler) UploadLabelAsset(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "Could not get user from context.")
		return
	}

	asset, err := h.FilesService.Upload(r, user.ID, "label", "label")
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("'%s' was uploaded successfully!", asset.FullFilename()))
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	err := h.FilesService.DownloadFile(w, r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, err.Error())
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	asset, err := h.FilesService.DeleteFile(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, err.Error())
		return
	}
	toastWarning(w, r, fmt.Sprintf("'%s' has been deleted.", asset.FullFilename()))
}
