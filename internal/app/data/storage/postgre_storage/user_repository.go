package postgrestorage

import (
	"red_test/internal/app/models"
)

type UserRepository struct {
	storage *PostgreStorage
}

func (r *UserRepository) CreateUser(u *models.User) error {
	_, err := r.storage.db.Exec(
		"INSERT INTO users (id, name, age, job) VALUES ($1, $2, $3, $4)",
		u.ID, u.Name, u.Age, u.Job,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(id string) (*models.User, error) {
	u := &models.User{}

	if err := r.storage.db.QueryRow(
		"SELECT id, name, age, job FROM users WHERE id = $1", id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Age,
		&u.Job,
	); err != nil {
		return nil, err
	}

	return u, nil
}
