package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, email string) (int, error)
	Get(ctx context.Context, id int) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
}

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, email string) (int, error) {
	var id int

	err := r.db.QueryRow(
		ctx,
		"INSERT INTO users(email) VALUES($1) RETURNING id",
		email,
	).Scan(&id)

	return id, err
}

func (r *Repo) Get(ctx context.Context, id int) (*User, error) {
	var u User

	err := r.db.QueryRow(
		ctx,
		"SELECT id, email FROM users WHERE id=$1",
		id,
	).Scan(&u.ID, &u.Email)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repo) GetAll(ctx context.Context) ([]User, error) {

	rows, err := r.db.Query(
		ctx,
		"SELECT id, email FROM users",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil

}
