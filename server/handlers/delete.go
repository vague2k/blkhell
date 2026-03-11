package handlers

import (
	"fmt"
	"net/http"

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

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	asset, err := h.FilesService.DeleteFile(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, err.Error())
		return
	}
	toastWarning(w, r, fmt.Sprintf("'%s' has been deleted.", asset.FullFilename()))
}
