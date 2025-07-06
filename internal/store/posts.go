package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id,omitempty"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version,omitempty"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post          Post `json:"post"`
	CommentsCount int  `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {

	query := `
		INSERT INTO posts (content, title, user_id, tags) 
		VALUES($1,$2,$3,$4) 
		RETURNING id, created_at, updated_at
	`

	// Though I added timeout to server config, and request alrdy has timeout
	// in r.Context() for write time duration,
	// here I'd like to make timeout shorter for DB
	// I can't override timeout passed but i can create another shorter timeout,
	// which will be fired first
	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
		SELECT p.id, p.content, p.title, p.user_id, p.tags, p.created_at, p.updated_at, p.version, u.username, u.avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1
	`

	var post Post

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
		&post.User.Username,
		&post.User.Avatar,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {

	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $5, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version, pq.Array(post.Tags)).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil

}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {

	query := `
		DELETE FROM posts
		WHERE posts.id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	// i dont care about rows affected because i expect that the row does exist
	// as i added middleware to fetch post first

	// rows, err := result.RowsAffected()
	// if err != nil {
	// 	return err
	// }

	// if rows == 0 {
	// 	return ErrNotFound
	// }

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {

	fmt.Println(fq)

	query := `
		SELECT p.id, p.title, p.content, p.tags, p.created_at, p.updated_at, p.user_id, u.username, u.avatar,   (
			SELECT COUNT(*) 
			FROM comments c 
			WHERE c.post_id = p.id
  		) AS comments_count
		FROM followers f
		JOIN posts p ON p.user_id = f.user_id OR p.user_id = $1
		JOIN users u ON u.id = p.user_id
		WHERE f.follower_id = $1 AND 
		(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
		(p.tags @> $5 OR $5 = '{}')
		GROUP BY p.id, p.title, p.content, p.tags, p.created_at, p.updated_at, p.user_id, u.username, u.avatar
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2
		OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}

	feed := make([]PostWithMetaData, 0, 10)

	for rows.Next() {
		var post PostWithMetaData
		err := rows.Scan(
			&post.Post.ID,
			&post.Post.Title,
			&post.Post.Content,
			pq.Array(&post.Post.Tags),
			&post.Post.CreatedAt,
			&post.Post.UpdatedAt,
			&post.Post.User.ID,
			&post.Post.User.Username,
			&post.Post.User.Avatar,
			&post.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, post)
	}

	return feed, nil
}

func (s *PostStore) GetUserPosts(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {

	fmt.Println(fq)

	query := `
		SELECT p.id, p.title, p.content, p.tags, p.created_at, p.updated_at, p.user_id, u.username, COUNT(c.id) AS comments_count
		FROM users u
		JOIN posts p ON p.user_id = u.id
		LEFT JOIN comments c ON c.post_id = p.id
		WHERE u.id = $1
		GROUP BY p.id, p.title, p.content, p.tags, p.created_at, p.updated_at, p.user_id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2
		OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, queryDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}

	feed := make([]PostWithMetaData, 0, 10)

	for rows.Next() {
		var post PostWithMetaData
		err := rows.Scan(
			&post.Post.ID,
			&post.Post.Title,
			&post.Post.Content,
			pq.Array(&post.Post.Tags),
			&post.Post.CreatedAt,
			&post.Post.UpdatedAt,
			&post.Post.User.ID,
			&post.Post.User.Username,
			&post.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, post)
	}

	return feed, nil
}
