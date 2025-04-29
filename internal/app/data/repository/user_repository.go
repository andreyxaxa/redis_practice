package repository

import "red_test/internal/app/models"

type UserRepository interface {
	CreateUser(*models.User) error
	GetUser(string) (*models.User, error)
}
