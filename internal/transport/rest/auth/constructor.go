package auth

import (
	"go.uber.org/fx"
	"notes_service/internal/services"
)

type Params struct {
	fx.In

	UserService services.UserService
	AuthService services.AuthService
}

type Handler struct {
	userService services.UserService
	authService services.AuthService
}

func NewHandler(p Params) *Handler {
	return &Handler{
		userService: p.UserService,
		authService: p.AuthService,
	}
}
