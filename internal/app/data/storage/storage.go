package storage

import "red_test/internal/app/data/repository"

type Storage interface {
	User() repository.UserRepository
}
