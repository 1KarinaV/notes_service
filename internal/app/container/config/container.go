package config

import (
	"go.uber.org/fx"
	"notes_service/internal/config"
)

func New() fx.Option {
	return fx.Provide(
		config.New,
	)
}
