package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vague2k/blkhell/server/database"
	serverErrors "github.com/vague2k/blkhell/server/errors"
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
	band, err := h.config.Database.GetBandByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band to display.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}
	pages.Band(&band).Render(r.Context(), w)
}

func (h *Handler) ReleasePage(w http.ResponseWriter, r *http.Request) {
	release, err := h.config.Database.GetReleaseByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get release to display.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}
	band, err := h.config.Database.GetBandByID(r.Context(), release.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			toastError(w, r, "Could not get band from release.")
			return
		}

		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	pages.Release(&band, &release).Render(r.Context(), w)
}

func (h *Handler) HXLabelAssetsImageGallery(w http.ResponseWriter, r *http.Request) {
	images, err := h.config.Database.GetLabelImageFiles(r.Context())
	if err != nil {
		toastError(w, r, "Could not get label files.")
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
	images, err := h.config.Database.GetBandImageFilesByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, "Could not get band images to display")
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
	images, err := h.config.Database.GetFileByPartialName(r.Context(), database.GetFileByPartialNameParams{
		Filename: "%" + input + "%",
		Ext:      "%" + input + "%",
	})
	if err != nil {
		toastError(w, r, "Search for asset image failed")
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
	stats, err := h.config.Database.GetDashboardStats(r.Context())
	if err != nil {
		toastError(w, r, "Could not get dashboard stats.")
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
		toastError(w, r, "Could not get user from context.")
		return
	}
	components.SidebarUserDropdown(user).Render(r.Context(), w)
}

func (h *Handler) HXSidebarBandsDropdown(w http.ResponseWriter, r *http.Request) {
	bands, err := h.config.Database.GetBands(r.Context())
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}
	components.SidebarBandsDropdown(bands).Render(r.Context(), w)
}

func (h *Handler) HXDashboardTable(w http.ResponseWriter, r *http.Request) {
	records, err := h.config.Database.GetDashboardBands(r.Context())
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
		return
	}

	components.DashboardTable(records).Render(r.Context(), w)
}

func (h *Handler) HXBandsReleaseTable(w http.ResponseWriter, r *http.Request) {
	releases, err := h.config.Database.GetReleasesByBand(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
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
	projects, err := h.config.Database.GetProjectsByBandID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, serverErrors.ErrDb.Error())
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
	err := h.FilesService.DownloadFile(w, r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		toastError(w, r, err.Error())
		return
	}
}
