package handlers

import (
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/server/services"
)

type Handler struct {
	AuthService      *services.AuthService
	BandsService     *services.BandsService
	FilesService     *services.FilesService
	DashboardService *services.DashboardService
	Middleware       *services.MiddlewareService
	DB               *database.Queries
}

func NewHandler(db *database.Queries) *Handler {
	return &Handler{
		AuthService:      services.NewAuthService(db),
		BandsService:     services.NewBandsService(db),
		FilesService:     services.NewFilesService(db),
		DashboardService: services.NewDashboardService(db),
		Middleware:       services.NewMiddlewareService(db),
		DB:               db,
	}
}
