package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.Auth.DestroySession(r)
	if err != nil {
		toastError(w, r, "500 Internal error: Could not destroy session.")
		return
	}

	w.Header().Set("HX-Redirect", "/login")
}

func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	image, err := h.DB.GetImageByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "500 Internal error: Could not delete image.")
		return
	}

	err = os.Remove(image.Path)
	if err != nil {
		toastError(w, r, "500 Internal error: Could not remove image from disk.")
		return
	}

	err = h.DB.DeleteImage(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "500 Internal error: Could not delete image.")
		return
	}

	toastWarning(w, r, fmt.Sprintf("'%s.%s' has been deleted.", image.Filename, image.Ext))
}
