package note

import (
	"go.uber.org/fx"
	"notes_service/internal/services"
)

type Params struct {
	fx.In

	AuthService services.AuthService
	NoteService services.NoteService
}

type Handler struct {
	authService services.AuthService
	noteService services.NoteService
}

func NewHandler(p Params) *Handler {
	return &Handler{
		authService: p.AuthService,
		noteService: p.NoteService,
	}
}
