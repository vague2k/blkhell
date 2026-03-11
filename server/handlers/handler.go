package handlers

import (
	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/server/services"
)

type Handler struct {
	AuthService      *services.AuthService
	FilesService     *services.FilesService
	DashboardService *services.DashboardService
	config           *config.Config
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{
		AuthService:      services.NewAuthService(config),
		FilesService:     services.NewFilesService(config),
		DashboardService: services.NewDashboardService(config),
		config:           config,
	}
}
