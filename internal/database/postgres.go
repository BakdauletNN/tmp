package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func Connect(dbURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Printf("unable to parse database url: %v", err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Printf("unable to connect to database: %v", err)
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("Unable to parse Database")
		pool.Close()
		return nil, err
	} else {
		log.Println("Connected to psql")
		return pool, nil

	}

}
