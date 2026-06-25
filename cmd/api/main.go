package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"grps-go-redis-psql/internal/app"
	"grps-go-redis-psql/internal/config"
	"grps-go-redis-psql/internal/db"
)

func main() {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.Load(
		configPath,
	)

	if err != nil {
		log.Fatal(err)
	}

	port := cfg.Server.Port

	dsn := os.Getenv("DB_DSN")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	pool, err := db.NewPostgresPool(ctx, dsn)
	if err != nil {
		log.Fatal("db init:", err)
	}
	defer pool.Close()

	application := app.New(pool, cfg)
	defer application.Close()

	log.Println("server started on port", port)

	if err := http.ListenAndServe(
		":"+port,
		application.Router(),
	); err != nil {
		log.Fatal(err)
	}
}
