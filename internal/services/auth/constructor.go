package auth

import (
	"go.uber.org/fx"
	"notes_service/internal/config"
	"notes_service/internal/services"
	"notes_service/pkg/web"
)

var _ services.AuthService = (*service)(nil)

type ServiceParams struct {
	fx.In

	Cfg   *config.Config
	Token *web.Token
}

type service struct {
	tokenLifetime int
	tkn           *web.Token
}

func NewService(p ServiceParams) services.AuthService {
	return &service{
		tokenLifetime: p.Cfg.TokenLifetime,
		tkn:           p.Token,
	}
}
