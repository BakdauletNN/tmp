package database

import (
	"context"
	"log"
	"tmp/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(db_URL string) (*pgxpool.Pool, error){
	var ctx context.Context = context.Background()
	var config *pgxpool.Pool
	var err error
	config, err = pgxpool.ParseConfig((db_URL))
	if err != nil{
		log.Printf("Unbale to parse Database url: %v", err)
		return nil, err
	}
	
}