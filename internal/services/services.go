package services

import (
	"accrual-loyalty-system-gophermart/internal/models"
	"accrual-loyalty-system-gophermart/internal/repository"
	"context"
)

type User interface {
	Create(ctx context.Context, login, password string) (*models.User, error)
	FindByLogin(ctx context.Context, login string) (*models.User, error)
}

type Services struct {
	User
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{
		User: NewUserService(repo.Users),
	}
}
