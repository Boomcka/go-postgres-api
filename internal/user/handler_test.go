package user

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

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
func TestHandler_CreateUser(t *testing.T) {
	repo := NewFakeRepo()
	h := NewHandler(repo)
	createUser(t, h, "test@test.com")
}
func TestHandler_GetUser(t *testing.T) {
	repo := NewFakeRepo()
	h := NewHandler(repo)
	var user = createUser(t, h, "test@test.com")
	getUser(t, h, user.ID)
}

func TestHandler_CreateAndGetUsers(t *testing.T) {
	repo := NewFakeRepo()
	h := NewHandler(repo)

	emails := []string{
		"first@test.com",
		"second@test.com",
		"third@test.com",
		"fourth@test.com",
		"fifth@test.com",
	}

	// ACT: create via HTTP, not repo
	for _, email := range emails {
		createUser(t, h, email)
	}

	// ACT: read via HTTP
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	h.Users(w, req)

	resp := w.Result()

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	// ASSERT
	if len(users) != len(emails) {
		t.Fatalf("expected %d users, got %d", len(emails), len(users))
	}
}
