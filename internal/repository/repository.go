package repository

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository/postgres"
)

type Users interface {
	Create(ctx context.Context, login, hashedPassword string) error
	Update(ctx context.Context, user *models.User) error
	Find(ctx context.Context, id int) (*models.User, error)
	FindByLogin(ctx context.Context, login string) (*models.User, error)
}

type Orders interface {
	Create(ctx context.Context, number string, user *models.User) error
	Update(ctx context.Context, order *models.Order) error
	FindByNumber(ctx context.Context, number string) (*models.Order, error)
	All(ctx context.Context, user *models.User) ([]*models.Order, error)
	AllPending(ctx context.Context) ([]*models.Order, error)
}

type Withdrawals interface {
	Create(ctx context.Context, withdrawal *models.Withdrawal) error
	Sum(ctx context.Context, user *models.User) (float64, error)
	All(ctx context.Context, user *models.User) ([]*models.Withdrawal, error)
}

type Repository struct {
	Users
	Orders
	Withdrawals
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		Users:       postgres.NewUserRepository(db),
		Orders:      postgres.NewOrderRepository(db),
		Withdrawals: postgres.NewWithdrawalRepository(db),
	}
}
