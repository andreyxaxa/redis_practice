package redisstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"red_test/internal/app/models"
)

type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) CreateUser(u *models.User) error {
	userKey := fmt.Sprintf("user:%s", u.ID)

	jsonUser, err := json.Marshal(u)
	if err != nil {
		return err
	}

	err = r.storage.client.Set(context.Background(), userKey, jsonUser, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(id string) (*models.User, error) {
	u := &models.User{}

	userKey := fmt.Sprintf("user:%s", id)

	uStr, err := r.storage.client.Get(context.Background(), userKey).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(uStr), u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
