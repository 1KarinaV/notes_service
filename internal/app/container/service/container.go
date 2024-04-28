package service

import (
	"go.uber.org/fx"
	"notes_service/internal/services/auth"
	"notes_service/internal/services/hash"
	"notes_service/internal/services/note"
	"notes_service/internal/services/user"
	"notes_service/pkg/web"
)

func New() fx.Option {
	return fx.Provide(
		web.NewToken,

		user.NewService,
		hash.NewService,
		auth.NewService,
		note.NewService,
	)
}
