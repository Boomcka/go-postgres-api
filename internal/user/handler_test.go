package user

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_CreateUser(t *testing.T) {
	repo := NewFakeRepo()
	h := NewHandler(repo)

	body := `{"email":"test@test.com"}`

	req := httptest.NewRequest(
		"POST",
		"/users",
		strings.NewReader(body),
	)

	w := httptest.NewRecorder()

	h.Users(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
func TestHandler_GetUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)

	defer cancel()
	repo := NewFakeRepo()
	h := NewHandler(repo)

	id, _ := repo.Create(ctx, "test@test.com")

	req := httptest.NewRequest(
		"GET",
		fmt.Sprintf("/users?id=%d", id),
		nil,
	)

	w := httptest.NewRecorder()

	h.Users(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHandler_CreateUser_InvalidEmail(t *testing.T) {
	repo := NewFakeRepo()
	h := NewHandler(repo)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "empty email",
			body:           `{"email":""}`,
			expectedStatus: 400,
		},
		{
			name:           "whitespace email",
			body:           `{"email":"   "}`,
			expectedStatus: 400,
		},
		{
			name:           "valid email",
			body:           `{"email":"test@test.com"}`,
			expectedStatus: 200,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			req := httptest.NewRequest(
				"POST",
				"/users",
				strings.NewReader(tc.body),
			)

			w := httptest.NewRecorder()

			h.Users(w, req)

			resp := w.Result()

			if resp.StatusCode != tc.expectedStatus {
				t.Fatalf(
					"%s: expected %d, got %d",
					tc.name,
					tc.expectedStatus,
					resp.StatusCode,
				)
			}
		})
	}
}
