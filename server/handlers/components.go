package handlers

import (
	"net/http"

	"github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/components"
)

func (h *Handler) HXSidebarUserDropdown(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "Could not get user from context.")
		return
	}
	components.SidebarUserDropdown(user).Render(r.Context(), w)
}

func (h *Handler) HXSidebarBandsDropdown(w http.ResponseWriter, r *http.Request) {
	bands, err := h.config.Database.GetActiveBands(r.Context())
	if err != nil {
		toastError(w, r, errors.ErrDb.Error())
		return
	}
	components.SidebarBandsDropdown(bands).Render(r.Context(), w)
}
