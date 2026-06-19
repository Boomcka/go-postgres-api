package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrEmailEmpty = errors.New("email is empty")
)

type Handler struct {
	repo UserRepository
}

type CreateUserRequest struct {
	Email string `json:"email"`
}

func NewHandler(repo UserRepository) *Handler {
	return &Handler{repo: repo}
}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")

	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, "bad id", 400)
		return
	}

	ctx := r.Context()

	user, err := h.repo.Get(ctx, id)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	if user == nil {
		http.Error(w, "not found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		h.create(w, r)

	case http.MethodGet:
		h.getAll(w, r)

	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (h *Handler) create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)
		return
	}
	if err := validateEmail(request.Email); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	ctx := r.Context()
	id, err := h.repo.Create(ctx, request.Email)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			http.Error(
				w,
				err.Error(),
				http.StatusConflict,
			)
			return
		}
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]any{
			"id": id},
	)
}

func validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return ErrEmailEmpty
	}
	return nil
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	users, err := h.repo.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	if len(users) == 0 {
		json.NewEncoder(w).Encode([]User{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
