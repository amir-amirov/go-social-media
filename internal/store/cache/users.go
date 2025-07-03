package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	redisDB *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {

	// { "user-42": 42 }
	cacheKey := fmt.Sprintf("user-%v", userID)

	data, err := s.redisDB.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Parse data from string to store.User
	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {

	if user == nil || user.ID == 0 {
		return errors.New("invalid user")
	}

	// { "user-42": 42 }
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.redisDB.SetEX(ctx, cacheKey, json, UserExpTime).Err()
}
