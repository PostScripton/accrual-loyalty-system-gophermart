package services

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/clients"
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

type Order interface {
	Create(ctx context.Context, number string, user *models.User) (*models.Order, error)
	FindByNumber(ctx context.Context, number string) (*models.Order, error)
	All(ctx context.Context, user *models.User) ([]*models.Order, error)
	RunPollingStatuses(ctx context.Context) error
}

type Luhn interface {
	Valid(number int) bool
}

type Services struct {
	User
	Auth
	Order
	Luhn
}

func NewServices(repo *repository.Repository, client *clients.AccrualSystemClient, JWTSecret string) *Services {
	return &Services{
		User:  NewUserService(repo.Users),
		Auth:  NewAuthService(repo.Users, JWTSecret),
		Order: NewOrderService(repo.Orders, client),
		Luhn:  &LuhnAlgo{},
	}
}
