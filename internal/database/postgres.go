package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"tmp/internal/logger"
)

func Connect(dbURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	log := logger.SetupLogger("local")

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Error("failed to parse database config", slog.String("db_url", dbURL), slog.String("error", err.Error()))
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("failed to create database pool", slog.String("error", err.Error()))
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Error("database ping failed", slog.String("error", err.Error()))
		pool.Close()
		return nil, err
	}

	log.Info("database connected")
	return pool, nil
}
