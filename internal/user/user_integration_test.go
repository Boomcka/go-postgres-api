package user

import (
	"encoding/json"
	"fmt"
	"grps-go-redis-psql/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIntegration_UserCreate(t *testing.T) {

	pool := testutil.SetupPostgres(t)
	repo := NewRepo(pool)
	handler := NewHandler(repo)
	createUser(t, handler, "test@test.com")

	req := httptest.NewRequest(
		"GET",
		"/users",
		nil,
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

func TestIntegration_DuplicateEmail(t *testing.T) {
	pool := testutil.SetupPostgres(t)
	repo := NewRepo(pool)
	handler := NewHandler(repo)
	createUser(t, handler, "test@emal.com")
	email := "test@emal.com"
	t.Helper()
	body := fmt.Sprintf(`{"email":"%s"}`, email)
	req := httptest.NewRequest(
		"POST",
		"/users",
		strings.NewReader(body),
	)

	w := httptest.NewRecorder()
	handler.Users(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf(
			"expected 409 got %d",
			resp.StatusCode,
		)
	}

}
