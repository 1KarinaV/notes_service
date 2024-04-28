package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"log"
	"notes_service/internal/config"
)

type Params struct {
	fx.In

	Cfg *config.Config
}

func NewPostgres(p Params) *sqlx.DB {
	db, err := sqlx.Connect(p.Cfg.DriverName, p.Cfg.DbConnect)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
