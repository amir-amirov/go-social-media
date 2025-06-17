package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID        int64     `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
}

type UserStore struct {
	db *sql.DB
}

func (s UserStore) Create(ctx context.Context, user *User) error {

	query := `
		INSERT INTO users(username, password, email) VALUES($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, &user.Username, &user.Password, &user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {

	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE id = $1
	`

	var user User

	if err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
