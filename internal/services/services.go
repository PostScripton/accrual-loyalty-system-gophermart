package services

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository"
)

type User interface {
	Create(ctx context.Context, login, password string) (*models.User, error)
	FindByLogin(ctx context.Context, login string) (*models.User, error)
}

type Auth interface {
	GetSecret() string
	LoginByUser(user *models.User) (string, error)
	Login(ctx context.Context, login, password string) (string, error)
}

type Services struct {
	User
	Auth
}

func NewServices(repo *repository.Repository, JWTSecret string) *Services {
	return &Services{
		User: NewUserService(repo.Users),
		Auth: NewAuthService(repo.Users, JWTSecret),
	}
}
