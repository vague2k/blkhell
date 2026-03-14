package handlers

import (
	"net/http"
	"net/url"

	"github.com/vague2k/blkhell/views/templui/toast"
)

type ToastFunc func(w http.ResponseWriter, r *http.Request, desc string) error

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

func toastWarning(w http.ResponseWriter, r *http.Request, desc string) error {
	return toast.Toast(toast.Props{
		Icon:        true,
		Title:       "Warning",
		Description: desc,
		Variant:     toast.VariantWarning,
		Position:    toast.PositionTopRight,
		Dismissible: true,
	}).Render(r.Context(), w)
}

func toastCookieError(w http.ResponseWriter, r *http.Request, cookieName string) {
	if msg := getToastCookieMsg(w, r, cookieName); msg != "" {
		toastError(w, r, msg)
	}
}

func toastCookieSuccess(w http.ResponseWriter, r *http.Request, cookieName string) {
	if msg := getToastCookieMsg(w, r, cookieName); msg != "" {
		toastSuccess(w, r, msg)
	}
}

func toastCookieWarning(w http.ResponseWriter, r *http.Request, cookieName string) {
	if msg := getToastCookieMsg(w, r, cookieName); msg != "" {
		toastWarning(w, r, msg)
	}
}

func setToastCookie(w http.ResponseWriter, cookieName string, toastMsg string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    url.QueryEscape(toastMsg),
		Path:     "/",
		MaxAge:   10,
		HttpOnly: true,
	})
}

func getToastCookieMsg(w http.ResponseWriter, r *http.Request, cookieName string) string {
	var msg string
	if c, err := r.Cookie(cookieName); err == nil {
		msg, err = url.QueryUnescape(c.Value)
		if err != nil {
			toastError(w, r, "can't unescape cookie value")
			return ""
		}

		// delete cookie after reading
		http.SetCookie(w, &http.Cookie{
			Name:     c.Name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
	}
	return msg
}
