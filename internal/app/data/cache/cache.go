package cache

import "red_test/internal/app/data/repository"

type Cache interface {
	User() repository.UserRepository
}
