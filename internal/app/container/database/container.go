package database

import (
	"go.uber.org/fx"
	"notes_service/internal/database/pg/note"
	"notes_service/internal/database/pg/user"
	"notes_service/pkg/postgres"
)

func New() fx.Option {
	return fx.Provide(
		// Data sources
		postgres.NewPostgres,

		user.NewRepository,
		note.NewRepository,
	)
}
