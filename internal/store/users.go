package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type User struct {
	ID        int64     `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"-"`
	IsActive  bool      `json:"is_active"`
	RoleID    int64     `json:"role_id"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) Compare(password string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(password))
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {

	query := `
		INSERT INTO users(
			username, 
			password, 
			email, 
			role_id
		) 
		VALUES($1, $2, $3, (SELECT id FROM roles WHERE name = 'user')) 
		RETURNING id, created_at;
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	// Sometimes I need this method to be a part of transaction so tx is introduced as optional parameter
	var err error
	if tx == nil {
		err = s.db.QueryRowContext(ctx, query, &user.Username, &user.Password.hash, &user.Email).Scan(&user.ID, &user.CreatedAt)
	} else {
		err = tx.QueryRowContext(ctx, query, &user.Username, &user.Password.hash, &user.Email).Scan(&user.ID, &user.CreatedAt)
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Constraint {
		case "users_email_key":
			return ErrDuplicateEmail
		case "users_username_key":
			return ErrDuplicateUsername
		}
	}

	return err
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	// transaction wrapper
	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// create the user invite
		if err := s.createUserInvitation(ctx, tx, token, user.ID, invitationExp); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {

	query := `
		SELECT id, username, email, password, created_at, role_id
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
		&user.Password.hash,
		&user.CreatedAt,
		&user.RoleID,
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

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {

	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE email = $1
	`

	var user User

	if err := s.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
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

// This is private method, so it cannot be used outside of the package
// Also it is not listed in interface list of methods for UserStore
func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userID int64, expiry time.Duration) error {

	query := `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, $3);
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(expiry))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// 1. find the user that this token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. update the user
		user.IsActive = true
		if err := s.Update(ctx, tx, user); err != nil {
			return err
		}

		// 3. delete the invitation
		if err := s.deleteInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil

	})
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {

	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND expiry > $2
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	var user User

	err := tx.QueryRowContext(ctx, query, token, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (s *UserStore) Update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users
		SET username = $2,
			email = $3,
			is_active = $4
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	var (
		result sql.Result
		err    error
	)

	if tx == nil {
		result, err = s.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.IsActive)
	} else {
		result, err = tx.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.IsActive)
	}

	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil || rows != 1 {
		return err
	}

	return nil

}

func (s *UserStore) deleteInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM user_invitations
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	return err
}
