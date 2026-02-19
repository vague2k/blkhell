package handlers

import (
	"github.com/vague2k/blkhell/server/auth"
	"github.com/vague2k/blkhell/server/database"
)

type Handler struct {
	Auth *auth.Service
	DB   *database.Queries
}

func NewHandler(authService *auth.Service, db *database.Queries) *Handler {
	return &Handler{
		Auth: authService,
		DB:   db,
	}
}
