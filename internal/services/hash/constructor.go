package hash

import (
	"go.uber.org/fx"
	"notes_service/internal/config"
	"notes_service/internal/services"
)

var _ services.HashService = (*service)(nil)

type Params struct {
	fx.In

	Cfg *config.Config
}

type service struct {
	hashCost int
}

func NewService(p Params) services.HashService {
	return &service{p.Cfg.HashCost}
}
