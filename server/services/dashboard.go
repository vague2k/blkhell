package services

import (
	"context"

	"github.com/vague2k/blkhell/config"
	serverErrors "github.com/vague2k/blkhell/server/errors"
	"github.com/vague2k/blkhell/views/templui/chart"
)

type DashboardService struct {
	config *config.Config
}

func NewDashboardService(cfg *config.Config) *DashboardService {
	return &DashboardService{config: cfg}
}

func (s *DashboardService) GetBandsFromPreviousYear(ctx context.Context) (*chart.Dataset, error) {
	bands, err := s.config.Database.GetBandsFromPreviousYear(ctx)
	if err != nil {
		return nil, serverErrors.ErrDb
	}
	bandCounts := make([]float64, 12)
	for _, band := range bands {
		m := int(band.CreatedAt.Month()) - 1
		bandCounts[m]++
	}

	return &chart.Dataset{
		Label: "Bands",
		Data:  bandCounts,
	}, nil
}

func (s *DashboardService) GetReleasesFromPreviousYear(ctx context.Context) (*chart.Dataset, error) {
	releases, err := s.config.Database.GetReleasesFromPreviousYear(ctx)
	if err != nil {
		return nil, serverErrors.ErrDb
	}
	releaseCounts := make([]float64, 12)
	for _, release := range releases {
		m := int(release.CreatedAt.Month()) - 1
		releaseCounts[m]++
	}

	return &chart.Dataset{
		Label: "Releases",
		Data:  releaseCounts,
	}, nil
}

func (s *DashboardService) GetProjectsFromPreviousYear(ctx context.Context) (*chart.Dataset, error) {
	projects, err := s.config.Database.GetProjectsFromPreviousYear(ctx)
	if err != nil {
		return nil, serverErrors.ErrDb
	}
	projectsCount := make([]float64, 12)
	for _, project := range projects {
		m := int(project.CreatedAt.Month()) - 1
		projectsCount[m]++
	}

	return &chart.Dataset{
		Label: "Projects",
		Data:  projectsCount,
	}, nil
}
