package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

type FollowStats struct {
	Followers int `json:"followers"`
	Following int `json:"following"`
}

func (s *FollowersStore) Follow(ctx context.Context, user_id int64, follower_id int64) error {

	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, follower_id) DO NOTHING;
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, user_id, follower_id)
	return err
}

func (s *FollowersStore) UnFollow(ctx context.Context, user_id int64, follower_id int64) error {

	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2;		
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, user_id, follower_id)
	return err
}

func (s *FollowersStore) Stats(ctx context.Context, user_id int64) (*FollowStats, error) {
	query := `
		SELECT COUNT(*) AS followers, 
		( 
			SELECT COUNT(*)
			FROM followers
			WHERE follower_id = $1
		) AS followings
		FROM followers
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	var stats FollowStats

	err := s.db.QueryRowContext(ctx, query, user_id).Scan(&stats.Followers, &stats.Following)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
