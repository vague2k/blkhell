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
	"github.com/vague2k/blkhell/views/templui/chart"
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

func (h *Handler) HXLabelAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.DB.GetLabelImageFiles(r.Context())
	if err != nil {
		toastError(w, r, "database error")
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)
	// render image count (hx-oob-swap)
	fmt.Fprintf(
		w,
		`<span id="label-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
		len(images),
	)
}

func (h *Handler) HXBandsAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.DB.GetBandImageFilesByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "database error")
		return
	}
	components.ImageGallery(images).Render(r.Context(), w)
	// render image count (hx-oob-swap)
	fmt.Fprintf(
		w,
		`<span id="band-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
		len(images),
	)
}

func (h *Handler) HXSearchLabelAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
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
		`<span id="label-assets-image-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d ASSETS</span>`,
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
		Title: "BANDS ON ROSTER",
		Count: strconv.FormatInt(stats.BandsCount, 10),
		Href:  fmt.Sprintf("/band/%s", stats.LatestBandID),
		Desc:  "View latest band",
		Name:  stats.LatestBandName.(string),
	}).Render(r.Context(), w)

	components.DashboardCard(components.DashboardCardProps{
		Title: "TOTAL RELEASES",
		Count: strconv.FormatInt(stats.ReleasesCount, 10),
		Href:  fmt.Sprintf("/release/%s", stats.LatestReleaseID),
		Desc:  "View latest release",
		Name:  stats.LatestReleaseTitle.(string),
	}).Render(r.Context(), w)

	components.DashboardCard(components.DashboardCardProps{
		Title: "LABEL ASSETS",
		Count: strconv.FormatInt(stats.LabelAssetsCount, 10),
		Href:  "/label-assets",
		Desc:  "Go to Blackheaven assets",
		Name:  "Assets",
	}).Render(r.Context(), w)

	components.DashboardCard(components.DashboardCardProps{
		Title: "PROJECTS",
		Count: strconv.FormatInt(stats.ProjectsCount, 10),
		Href:  fmt.Sprintf("/project/%s", stats.LatestProjectID),
		Desc:  "Go to latest project",
		Name:  stats.LatestProjectName.(string),
	}).Render(r.Context(), w)
}

func (h *Handler) HXDashboardChart(w http.ResponseWriter, r *http.Request) {
	bandsDataSet, err := h.DashboardService.GetBandsFromPreviousYear(r.Context())
	if err != nil {
		toastError(w, r, err.Error())
		return
	}
	releasesDataSet, err := h.DashboardService.GetReleasesFromPreviousYear(r.Context())
	if err != nil {
		toastError(w, r, err.Error())
		return
	}
	projectsDataSet, err := h.DashboardService.GetProjectsFromPreviousYear(r.Context())
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	chart.Chart(chart.Props{
		Class:       "h-full",
		Variant:     chart.VariantLine,
		ShowYGrid:   true,
		ShowYLabels: true,
		ShowXLabels: true,
		ShowLegend:  true,
		Data: chart.Data{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Oct", "Sep", "Nov", "Dec"},
			Datasets: []chart.Dataset{
				*bandsDataSet,
				*releasesDataSet,
				*projectsDataSet,
			},
		},
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

func (h *Handler) HXBandsReleaseTable(w http.ResponseWriter, r *http.Request) {
	releases, err := h.DB.GetReleasesByBand(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "database error")
		return
	}

	components.BandReleasesTable(releases).Render(r.Context(), w)
	// render release count (hx-oob-swap)
	fmt.Fprintf(
		w,
		`<span id="band-releases-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d RELEASES</span>`,
		len(releases),
	)
}

func (h *Handler) HXBandProjectsTable(w http.ResponseWriter, r *http.Request) {
	projects, err := h.DB.GetProjectsByBandID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "database error")
		return
	}

	components.BandProjectsTable(projects).Render(r.Context(), w)
	// render release count (hx-oob-swap)
	if len(projects) > 0 {
		fmt.Fprintf(
			w,
			`<span id="band-projects-count" hx-swap-oob="true" class="font-light text-muted-foreground text-sm">%d PROJECTS</span>`,
			len(projects),
		)
	}
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	asset, err := h.DB.GetFileByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(os.Getenv("UPLOADS_DIR") + asset.Path)
	if err != nil {
		http.Error(w, "download failed", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, asset.FullFilename()))

	io.Copy(w, file)
}
