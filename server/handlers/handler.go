package handlers

import (
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/server/services"
)

type Handler struct {
	AuthService      *services.AuthService
	FilesService     *services.FilesService
	DashboardService *services.DashboardService
	Middleware       *services.MiddlewareService
	DB               *database.Queries
}

func NewHandler(db *database.Queries) *Handler {
	return &Handler{
		AuthService:      services.NewAuthService(db),
		FilesService:     services.NewFilesService(db),
		DashboardService: services.NewDashboardService(db),
		Middleware:       services.NewMiddlewareService(db),
		DB:               db,
	}
}
