package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")

	queryDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetaData, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Activate(context.Context, string) error
		Update(context.Context, *sql.Tx, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByID(context.Context, int64) (*Comment, error)
		GetByPostID(context.Context, int64) ([]Comment, error)
		Update(context.Context, *Comment) error
		Delete(context.Context, int64) error
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		UnFollow(context.Context, int64, int64) error
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
		GetByID(context.Context, int64) (*Role, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowersStore{db},
		Roles:     &RoleStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
