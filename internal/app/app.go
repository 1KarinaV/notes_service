package app

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"notes_service/internal/app/container/config"
	"notes_service/internal/app/container/database"
	"notes_service/internal/app/container/service"
	"notes_service/internal/app/container/transport"
	transp "notes_service/internal/transport/rest"
	"time"
)

type App struct {
	fxOptions  fx.Option
	fxApp      *fx.App
	httpServer *transp.Server
}

func New() (*App, error) {
	app := new(App)

	app.FxProvides(
		service.New,
		config.New,
		transport.New,
		database.New,
	)

	return app, nil
}

func (a *App) FxProvides(ff ...func() fx.Option) {
	options := make([]fx.Option, len(ff))
	for rs, fn := range ff {
		options[rs] = fn()
	}
	a.fxOptions = fx.Options(options...)
}

func (a *App) Init() error {
	a.fxOptions = fx.Options(
		fx.StopTimeout(time.Minute*10),
		a.fxOptions,
		fx.NopLogger,

		fx.Populate(&a.httpServer),
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: a.onStart,
				OnStop:  a.onStop,
			})
		}),
	)

	a.fxApp = fx.New(a.fxOptions)

	return nil
}

func (a *App) onStart(ctx context.Context) error {
	zap.L().Info("Starting")
	if err := a.httpServer.Start(); err != nil {
		return err
	}
	zap.L().Info("Started")
	return nil
}

func (a *App) onStop(ctx context.Context) error {
	zap.L().Info("Stopping")
	if err := a.httpServer.Stop(); err != nil {
		return err
	}
	zap.L().Info("Stopped")
	return nil
}

func (a *App) Run() error {
	if err := a.fxApp.Err(); err != nil {
		return err
	}

	a.fxApp.Run()
	return nil
}
