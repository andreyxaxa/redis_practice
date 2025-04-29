package rediscache

import (
	"red_test/internal/app/data/repository"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client         *redis.Client
	userRepository repository.UserRepository
}

func New(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) User() repository.UserRepository {
	if c.userRepository != nil {
		return c.userRepository
	}

	c.userRepository = &UserRepository{
		redis: c,
	}
	return c.userRepository
}
