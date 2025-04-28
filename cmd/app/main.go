package main

import (
	"context"
	"io-load-api/internal/app"
	"io-load-api/internal/config"
	"io-load-api/internal/metrics"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	logger.Info("Starting app...")
	logger.Info("Start metrics")
	metrics.RegisterMetrics()
	metrics.StartMetricsServer(cfg)
	logger.Info("Initializing app...")
	application, err := app.New(logger, cfg)
	if err != nil {
		log.Fatalf("error initializing application %s", err)
	}

	// Running application
	logger.Info("Running app...")
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- application.MustRun()
	}()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	logger.Info("Application started")

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
