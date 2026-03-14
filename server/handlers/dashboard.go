package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/components"
	"github.com/vague2k/blkhell/views/pages"
	"github.com/vague2k/blkhell/views/templui/chart"
)

func (h *Handler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	toastCookieWarning(w, r, "deleted_band_toast_message")
	pages.Dashboard().Render(r.Context(), w)
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

func (h *Handler) HXDashboardTable(w http.ResponseWriter, r *http.Request) {
	records, err := h.config.Database.GetDashboardBands(r.Context())
	if err != nil {
		toastError(w, r, errors.ErrDb.Error())
		return
	}

	components.DashboardTable(records).Render(r.Context(), w)
}
