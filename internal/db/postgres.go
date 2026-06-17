package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() *pgxpool.Pool {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("failed to create pool:", err)
	}

	// проверка подключения
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to ping db:", err)
	}

	return pool
}
