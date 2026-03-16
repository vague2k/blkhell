package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/common"
	"github.com/vague2k/blkhell/server/database"
	serverErrors "github.com/vague2k/blkhell/server/errors"
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

	count := len(images)
	if count <= 0 {
		fmt.Fprint(
			w,
			`<span id="label-assets-image-count" hx-swap-oob="true"></span>`,
		)
	} else {
		fmt.Fprintf(
			w,
			`<span id="label-assets-image-count" hx-swap-oob="true" class="text-muted-foreground text-xs">%d ASSETS</span>`,
			count,
		)
	}
}

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

func (h *Handler) HXReleaseAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.config.Database.GetReleaseImageFilesByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "Could not get releases images to display")
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)

	count := len(images)
	if count <= 0 {
		fmt.Fprint(
			w,
			`<span id="release-assets-image-count" hx-swap-oob="true"></span>`,
		)
	} else {
		fmt.Fprintf(
			w,
			`<span id="release-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
			count,
		)
	}
}

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

func (h *Handler) UploadLabelAsset(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "Could not get user from context.")
		return
	}

	asset, err := h.FilesService.Upload(w, r, user.ID, "label", common.FileOwnerTypeLabel)
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

func (h *Handler) EditFile(w http.ResponseWriter, r *http.Request) {
	renamed := r.FormValue("filename")
	if renamed == "" {
		toastError(w, r, "At least 1 field has to be filled.")
		return
	}

	file, err := h.config.Database.GetFileByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get file to edit.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	params := database.UpdateFileParams{
		ID:       file.ID,
		Filename: file.Filename,
	}

	if renamed != "" {
		params.Filename = renamed
	}

	_, err = h.config.Database.UpdateFile(r.Context(), params)
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	toastSuccess(w, r, "Your changes have been saved!")
}
