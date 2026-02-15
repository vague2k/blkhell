package handlers

import "github.com/vague2k/blkhell/server/auth"

type Handler struct {
	auth *auth.Service
}

func NewHandler(authService *auth.Service) *Handler {
	return &Handler{
		auth: authService,
	}
}
