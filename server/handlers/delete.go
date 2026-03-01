package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.AuthService.DestroySession(w, r)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	w.Header().Set("HX-Redirect", "/login")
}

func (h *Handler) DeleteLabelAsset(w http.ResponseWriter, r *http.Request) {
	asset, err := h.DB.DeleteFile(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "500 Internal error: Could not delete image.")
		return
	}

	// TODO: make more robust, perhaps put os call before db call ?
	uploadsDir := os.Getenv("UPLOADS_DIR")
	err = os.Remove(uploadsDir + asset.Path)
	if err != nil {
		toastError(w, r, "500 Internal error: Could not remove image from disk.")
		return
	}

	toastWarning(w, r, fmt.Sprintf("'%s.%s' has been deleted.", asset.Filename, asset.Ext))
}
