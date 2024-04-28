package web

import (
	"github.com/go-chi/jwtauth"
	"go.uber.org/fx"
	"notes_service/internal/config"
)

type Params struct {
	fx.In

	Cfg *config.Config
}

type Token struct {
	*jwtauth.JWTAuth
}

func NewToken(p Params) *Token {
	return &Token{
		jwtauth.New("HS256", []byte(p.Cfg.SecretKey), nil),
	}
}
