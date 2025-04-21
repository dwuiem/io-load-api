package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"test-workmate/internal/config"
	"time"
)

type Store struct {
	db *pgxpool.Pool
}

func New(cfg *config.Config) (Store, error) {
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

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		err := pool.Ping(context.Background())
		if err == nil {
			return Store{db: pool}, nil
		}
		time.Sleep(time.Second)
	}
	return Store{}, errors.New("failed to connect to postgres")
}
