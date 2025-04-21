package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"test-workmate/internal/app"
	"test-workmate/internal/config"
	"test-workmate/internal/metrics"
	"time"
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	metrics.RegisterMetrics()
	metrics.StartMetricsServer(cfg)
	application, err := app.New(logger, cfg)
	if err != nil {
		log.Fatal("error initializing application", err)
	}

	// Running application
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- application.MustRun()
	}()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	// Graceful shutdown
	select {
	case err := <-serverErrors:
		logger.Info("Application exited unexpectedly", slog.String("error", err.Error()))
	case sig := <-stopSignal:
		logger.Info("Stopping application", slog.String("signal", sig.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := application.Stop(ctx); err != nil {
			logger.Info("Application exited with error", slog.String("error", err.Error()))
		} else {
			logger.Info("Application exited gracefully")
		}
	}
}
