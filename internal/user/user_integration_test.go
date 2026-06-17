package user

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestIntegration_UserCreate(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewRepo(pool)
	handler := NewHandler(repo)
	req := httptest.NewRequest(
		"POST",
		"/users",
		strings.NewReader(`{"email":"test@test.com"}`),
	)
	w := httptest.NewRecorder()
	handler.Users(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf(
			"%s: expected %d, got %d",
			"valid email",
			200,
			resp.StatusCode,
		)
	}

	req = httptest.NewRequest(
		"GET",
		"/users",
		nil,
	)

	w = httptest.NewRecorder()

	handler.Users(w, req)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf(
			"%s: expected %d, got %d",
			"valid email",
			200,
			resp.StatusCode,
		)
	}
	var users []User
	err := json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if users[0].Email != "test@test.com" {
		t.Fatalf(
			"%s: expected %d, got %d",
			"valid email",
			200,
			resp.StatusCode,
		)
	}
}

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)

	t.Cleanup(cancel)

	container, err := postgres.Run(
		ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").
				WithStartupTimeout(60*time.Second),
		),
	)

	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	connStr, err := container.ConnectionString(
		ctx,
		"sslmode=disable",
	)

	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)

	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping postgres: %v", err)
	}

	migration, err := os.ReadFile(
		"../migrations/001_users.sql",
	)

	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	_, err = pool.Exec(
		ctx,
		string(migration),
	)

	if err != nil {
		t.Fatalf("run migration: %v", err)
	}

	return pool
}
