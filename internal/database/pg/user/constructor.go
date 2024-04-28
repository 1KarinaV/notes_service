package user

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"notes_service/internal/database"
)

var _ database.UserRp = (*repository)(nil)

type Params struct {
	fx.In

	DB *sqlx.DB
}

type repository struct {
	*sqlx.DB
}

func NewRepository(p Params) database.UserRp {
	return &repository{
		p.DB,
	}
}
