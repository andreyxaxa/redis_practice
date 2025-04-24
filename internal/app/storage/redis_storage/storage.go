package redisstorage

import (
	"red_test/internal/app/storage"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client         *redis.Client
	userRepository storage.UserRepository
}

func New(addr string) *Storage {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Storage{
		client: rdb,
	}
}

func (s *Storage) User() storage.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		storage: s,
	}
	return s.userRepository
}
