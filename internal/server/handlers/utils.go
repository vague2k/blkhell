package handlers

import (
	"net/http"

	"github.com/vague2k/blkhell/views/templui/toast"
)

func toastError(w http.ResponseWriter, r *http.Request, desc string) error {
	return toast.Toast(toast.Props{
		Icon:        true,
		Title:       "Error",
		Description: desc,
		Variant:     toast.VariantError,
		Position:    toast.PositionTopRight,
		Dismissible: true,
	}).Render(r.Context(), w)
}

func toastSuccess(w http.ResponseWriter, r *http.Request, desc string) error {
	return toast.Toast(toast.Props{
		Icon:        true,
		Title:       "Success",
		Description: desc,
		Variant:     toast.VariantSuccess,
		Position:    toast.PositionTopRight,
		Dismissible: true,
	}).Render(r.Context(), w)
}
