package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/views/components"
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
	pages.Dashboard().Render(r.Context(), w)
}

func (h *Handler) LabelAssetsPage(w http.ResponseWriter, r *http.Request) {
	pages.LabelAssets().Render(r.Context(), w)
}

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	pages.Settings().Render(r.Context(), w)
}

func (h *Handler) BandPage(w http.ResponseWriter, r *http.Request) {
	band, err := h.DB.GetBandByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "database error")
		return
	}
	pages.Band(&band).Render(r.Context(), w)
}

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
		`<span id="dashboard-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
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
		`<span id="dashboard-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
		len(images),
	)
}

func (h *Handler) HXDashboardCards(w http.ResponseWriter, r *http.Request) {
	stats, err := h.DB.GetDashboardStats(r.Context())
	if err != nil {
		toastError(w, r, "database error")
		return
	}

	components.DashboardCard(components.DashboardCardProps{
		Title: "LABEL ASSETS",
		Count: strconv.FormatInt(stats.LabelAssets, 10),
		Icon:  icon.FileBox(icon.Props{Size: 20}),
	}).Render(r.Context(), w)

	components.DashboardCard(components.DashboardCardProps{
		Title: "SIGNED BANDS",
		Count: strconv.FormatInt(stats.Bands, 10),
		Icon:  icon.Users(icon.Props{Size: 20}),
	}).Render(r.Context(), w)

	components.DashboardCard(components.DashboardCardProps{
		Title: "TOTAL RELEASES",
		Count: strconv.FormatInt(stats.Releases, 10),
		Icon:  icon.DiscAlbum(icon.Props{Size: 20}),
	}).Render(r.Context(), w)

	components.DashboardCard(components.DashboardCardProps{
		Title: "PROJECTS",
		Count: strconv.FormatInt(stats.Projects, 10),
		Icon:  icon.FolderArchive(icon.Props{Size: 20}),
	}).Render(r.Context(), w)
}

func (h *Handler) HXSidebarUserDropdown(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user in context", http.StatusInternalServerError)
		return
	}
	components.SidebarUserDropdown(user).Render(r.Context(), w)
}

func (h *Handler) HXSidebarBandsDropdown(w http.ResponseWriter, r *http.Request) {
	bands, err := h.DB.GetBands(r.Context())
	if err != nil {
		toastError(w, r, "database error")
		return
	}
	components.SidebarBandsDropdown(bands).Render(r.Context(), w)
}

func (h *Handler) HXDashboardTable(w http.ResponseWriter, r *http.Request) {
	records, err := h.DB.GetDashboardBands(r.Context())
	if err != nil {
		toastError(w, r, "database error")
		return
	}

	components.DashboardTable(records).Render(r.Context(), w)
}

func (h *Handler) DownloadLabelAsset(w http.ResponseWriter, r *http.Request) {
	asset, err := h.DB.GetFileByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(os.Getenv("UPLOADS_DIR") + asset.Path)
	if err != nil {
		toastError(w, r, "Couldn't open file to download")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.%s"`, asset.Filename, asset.Ext))

	io.Copy(w, file)
}
