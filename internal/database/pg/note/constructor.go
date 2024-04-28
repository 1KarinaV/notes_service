package note

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"notes_service/internal/database"
)

var _ database.NoteRp = (*repository)(nil)

type Params struct {
	fx.In

	DB *sqlx.DB
}

type repository struct {
	*sqlx.DB
}

func NewRepository(p Params) database.NoteRp {
	return &repository{
		p.DB,
	}
}
