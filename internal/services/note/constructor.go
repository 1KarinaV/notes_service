package note

import (
	"go.uber.org/fx"
	"notes_service/internal/database"
	"notes_service/internal/services"
)

type Params struct {
	fx.In

	NoteRp  database.NoteRp
	UserSrv services.UserService
}

type service struct {
	noteRp  database.NoteRp
	userSrv services.UserService
}

func NewService(p Params) services.NoteService {
	return &service{
		noteRp:  p.NoteRp,
		userSrv: p.UserSrv,
	}
}
