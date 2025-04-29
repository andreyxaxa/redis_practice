package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"red_test/internal/app/models"
	"time"
)

type UserRepository struct {
	redis *RedisCache
}

// SET
func (r *UserRepository) CreateUser(u *models.User) error {
	userKey := fmt.Sprintf("user:%s", u.ID)

	userValue, err := json.Marshal(u)
	if err != nil {
		return err
	}

	err = r.redis.client.Set(context.Background(), userKey, userValue, time.Hour*24).Err()
	if err != nil {
		return err
	}

	return nil
}

// GET
func (r *UserRepository) GetUser(id string) (*models.User, error) {
	u := &models.User{}

	userKey := fmt.Sprintf("user:%s", id)

	userValue, err := r.redis.client.Get(context.Background(), userKey).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(userValue), u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// DEL
