package postgrestorage

import (
	"database/sql"
	"red_test/internal/app/data/repository"
)

type PostgreStorage struct {
	db             *sql.DB
	userRepository repository.UserRepository
}

func New(db *sql.DB) *PostgreStorage {
	return &PostgreStorage{
		db: db,
	}
}

func (s *PostgreStorage) User() repository.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		storage: s,
	}

	return s.userRepository
}
