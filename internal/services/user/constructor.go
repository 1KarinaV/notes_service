package user

import (
	"go.uber.org/fx"
	"notes_service/internal/database"
	"notes_service/internal/services"
)

var _ services.UserService = (*service)(nil)

type Params struct {
	fx.In

	UserRp      database.UserRp
	HashService services.HashService
}

type service struct {
	userRp      database.UserRp
	hashService services.HashService
}

func NewService(p Params) services.UserService {
	return &service{
		userRp:      p.UserRp,
		hashService: p.HashService,
	}
}
