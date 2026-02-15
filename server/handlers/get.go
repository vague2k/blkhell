package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/vague2k/blkhell/views/pages"
)

func (h *Handler) LoginPage() http.Handler {
	return templ.Handler(pages.Login())
}

func (h *Handler) DashboardPage() http.Handler {
	return templ.Handler(pages.Dashboard())
}
