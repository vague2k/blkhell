package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/views/components"
	"github.com/vague2k/blkhell/views/layouts"
	"github.com/vague2k/blkhell/views/pages"
	"github.com/vague2k/blkhell/views/templui/icon"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	// redirects user if non expired session is already found for user from request
	if ok && user != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	pages.Login().Render(r.Context(), w)
}

func (h *Handler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user in context", http.StatusInternalServerError)
		return
	}
	layouts.BaseSidebarLayout(user, pages.Dashboard()).Render(r.Context(), w)
}

func (h *Handler) LabelAssetsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user in context", http.StatusInternalServerError)
		return
	}
	layouts.BaseSidebarLayout(user, pages.LabelAssets()).Render(r.Context(), w)
}

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user in context", http.StatusInternalServerError)
		return
	}
	layouts.BaseSidebarLayout(user, pages.Settings()).Render(r.Context(), w)
}

// func (h *Handler) BandsPage(w http.ResponseWriter, r *http.Request) {
// 	user, ok := h.AuthService.UserFromContext(r.Context())
// 	if !ok {
// 		http.Error(w, "missing user in context", http.StatusInternalServerError)
// 		return
// 	}
// 	bands, ok := h.BandsService.BandsFromContext(r.Context())
// 	if !ok {
// 		http.Error(w, "missing bands in context", http.StatusInternalServerError)
// 		return
// 	}
// 	layouts.BaseSidebarLayout(user, pages.BandsPage(bands)).Render(r.Context(), w)
// }

func (h *Handler) HXImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.DB.GetFiles(r.Context())
	if err != nil {
		http.Error(w, "Couldn't get user from request", http.StatusInternalServerError)
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)
	// render image count (hx-oob-swap)
	fmt.Fprintf(
		w,
		`<span id="dashboard-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d IMAGES</span>`,
		len(images),
	)
}

func (h *Handler) HXSearchImageGallery(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("q")
	images, err := h.DB.GetFileByPartialName(r.Context(), database.GetFileByPartialNameParams{
		Filename: "%" + input + "%",
		Ext:      "%" + input + "%",
	})
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	components.ImageGallery(images).Render(r.Context(), w)
	// render image count (hx-oob-swap)
	fmt.Fprintf(
		w,
		`<span id="dashboard-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d IMAGES</span>`,
		len(images),
	)
}

func (h *Handler) HXDashboardCards(w http.ResponseWriter, r *http.Request) {
	bands, err := h.DB.GetBands(r.Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		toastError(w, r, "database error")
		return
	}
	releases, err := h.DB.GetReleases(r.Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		toastError(w, r, "database error")
		return
	}
	projects, err := h.DB.GetProjects(r.Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		toastError(w, r, "database error")
		return
	}

	components.DashboardCard(components.DashboardCardProps{
		Title: "SIGNED BANDS",
		Count: fmt.Sprintf("%d", len(bands)),
		Icon:  icon.Users(icon.Props{Size: 20}),
	}).Render(r.Context(), w)
	components.DashboardCard(components.DashboardCardProps{
		Title: "TOTAL RELEASES",
		Count: fmt.Sprintf("%d", len(releases)),
		Icon:  icon.DiscAlbum(icon.Props{Size: 20}),
	}).Render(r.Context(), w)
	components.DashboardCard(components.DashboardCardProps{
		Title: "PROJECTS",
		Count: fmt.Sprintf("%d", len(projects)),
		Icon:  icon.FolderArchive(icon.Props{Size: 20}),
	}).Render(r.Context(), w)
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	image, err := h.DB.GetFileByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(image.Path)
	if err != nil {
		toastError(w, r, "Couldn't open file to download")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.%s"`, image.Filename, image.Ext))

	io.Copy(w, file)
}
