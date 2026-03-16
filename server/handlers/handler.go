package handlers

import (
	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/server/services"
)

type Handler struct {
	authService      *services.AuthService
	filesService     *services.FilesService
	dashboardService *services.DashboardService
	config           *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		authService:      services.NewAuthService(cfg),
		filesService:     services.NewFilesService(cfg),
		dashboardService: services.NewDashboardService(cfg),
		config:           cfg,
	}
}
