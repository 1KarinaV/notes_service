package transport

import (
	"go.uber.org/fx"
	"notes_service/internal/transport/rest"
	"notes_service/internal/transport/rest/auth"
	"notes_service/internal/transport/rest/note"
)

func New() fx.Option {
	return fx.Provide(
		// Rest
		rest.New,
		auth.NewHandler,
		note.NewHandler,

		// Rpc
	)
}
