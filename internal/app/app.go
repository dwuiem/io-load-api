package app

import (
	"context"
	"log/slog"
	"net/http"
	"test-workmate/internal/config"
	"test-workmate/internal/service"
	"test-workmate/internal/store/postgres"
	"test-workmate/internal/transport/http/handler"
)

// App contains HTTP Server, initializes store, handler and services and runs server
type App struct {
	HTTPServer *http.Server
	log        *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	store, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}
	taskStore := postgres.NewTaskStore(store)
	services := service.NewTaskService(log, taskStore)
	handlers := handler.New(log, services)
	return &App{
		HTTPServer: &http.Server{
			Addr:    cfg.HTTPServer.Addr,
			Handler: handlers.InitRoutes(),
		},
		log: log,
	}, nil
}
func (app *App) MustRun() error {
	app.log.Info("Running HTTP server")
	return app.HTTPServer.ListenAndServe()
}

func (app *App) Stop(ctx context.Context) error {
	app.log.Info("Stopping HTTP server")
	return app.HTTPServer.Shutdown(ctx)
}
