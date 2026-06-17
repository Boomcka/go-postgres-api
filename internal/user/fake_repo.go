package user

import (
	"context"
	"sync"
)

type FakeRepo struct {
	mu           sync.Mutex
	users        map[int]User
	next         int
	createrCalls int
}

func NewFakeRepo() *FakeRepo {
	return &FakeRepo{
		users: make(map[int]User),
		next:  1,
	}
}
func (r *FakeRepo) Create(ctx context.Context, email string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.createrCalls++
	id := r.next
	r.next++

	r.users[id] = User{
		ID:    id,
		Email: email,
	}

	return id, nil
}
func (r *FakeRepo) Get(ctx context.Context, id int) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[id]
	if !ok {
		return nil, nil
	}

	return &u, nil
}

func (r *FakeRepo) GetAll(ctx context.Context) ([]User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u := make([]User, 0, len(r.users))
	for _, user := range r.users {
		u = append(u, user)
	}
	return u, nil
}

func (r *FakeRepo) GetCreaterCalls() (calls int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.createrCalls
}
