package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	log.Println("Connected to Database")
	return pool, nil
}
