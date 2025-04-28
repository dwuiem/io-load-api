package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io-load-api/internal/config"
	"log/slog"
	"time"
)

type Store struct {
	db *pgxpool.Pool
}

func New(log *slog.Logger, cfg *config.Config) (Store, error) {
	pgxConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.PostgresDB.Username,
			cfg.PostgresDB.Password,
			cfg.PostgresDB.Host,
			cfg.PostgresDB.Port,
			cfg.PostgresDB.DBName,
		),
	)

	pgxConfig.MaxConns = int32(cfg.PostgresDB.MaxConns)
	pgxConfig.MinConns = int32(cfg.PostgresDB.MinConns)
	pgxConfig.MaxConnIdleTime = cfg.PostgresDB.MaxConnIdleTime
	pgxConfig.HealthCheckPeriod = cfg.PostgresDB.HealthCheckPeriod

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Initializing pgx pool")
	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return Store{}, errors.New("failed to create pgx pool")
	}
	for i := 0; i < 5; i++ {
		log.Info("Ping database")
		pingCtx, pingCancel := context.WithTimeout(context.Background(), time.Second)
		err := pool.Ping(pingCtx)
		pingCancel()
		if err == nil {
			return Store{db: pool}, nil
		}
	}

	pool.Close()
	return Store{}, errors.New("failed to connect to postgres")
}
