package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"notes_service/internal/config"
	"notes_service/internal/transport/rest/auth"
	"notes_service/internal/transport/rest/note"
	"notes_service/pkg/web"
	"os"
	"os/signal"
	"time"
)

type Params struct {
	fx.In
	Cfg *config.Config

	AuthHandler *auth.Handler
	NoteHandler *note.Handler

	Token *web.Token
}

type Server struct {
	*http.Server
}

func New(p Params) (*Server, error) {

	r := chi.NewRouter()

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "origin", "x-requested-with"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Middlewares
	r.Use(middleware.URLFormat)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(jwtauth.Verifier(p.Token.JWTAuth))

	r.Group(func(r chi.Router) {
		r.Use(p.AuthHandler.AuthCtx)
		r.Post("/register", p.AuthHandler.Register)
		r.Post("/login", p.AuthHandler.Login)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator)
		r.Get("/logout", p.AuthHandler.Logout)
	})
	r.Route("/note", func(r chi.Router) {

		r.Get("/list", p.NoteHandler.ListNotes)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Authenticator)
			r.Post("/create", p.NoteHandler.CreateNote)
			r.Delete("/delete", p.NoteHandler.DeleteNote)
			r.Put("/update", p.NoteHandler.UpdateNote)
		})
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", p.Cfg.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		srv,
	}, nil
}

func (s *Server) Start() error {
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("Could not listen on", zap.String("addr", s.Addr), zap.Error(err))
		}
	}()

	zap.L().Info("Server is ready to handle requests", zap.String("addr", s.Addr))
	return nil
}

func (s *Server) Stop() error {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		zap.L().Fatal("Could not gracefully shutdown the server", zap.Error(err))
	}

	zap.L().Info("Server stopped")
	return nil
}
