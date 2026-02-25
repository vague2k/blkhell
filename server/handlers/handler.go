package handlers

import (
	"github.com/vague2k/blkhell/server/database"
	"github.com/vague2k/blkhell/server/services"
)

type Handler struct {
	AuthService *services.AuthService
	DB          *database.Queries
}

func NewHandler(authService *services.AuthService, db *database.Queries) *Handler {
	return &Handler{
		AuthService: authService,
		DB:          db,
	}
}
