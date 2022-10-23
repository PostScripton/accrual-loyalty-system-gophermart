package repository

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Users interface {
	Create(ctx context.Context, login, hashedPassword string) error
	FindByLogin(ctx context.Context, login string) (*models.User, error)
}

type Orders interface {
	Create(ctx context.Context, number string, user *models.User) error
	Update(ctx context.Context, order *models.Order) error
	FindByNumber(ctx context.Context, number string, user *models.User) (*models.Order, error)
	All(ctx context.Context, user *models.User) ([]*models.Order, error)
	AllPending(ctx context.Context) ([]*models.Order, error)
}

type Repository struct {
	Users
	Orders
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Users:  postgres.NewUserRepository(db),
		Orders: postgres.NewOrderRepository(db),
	}
}
