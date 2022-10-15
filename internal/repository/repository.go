package repository

import (
	"accrual-loyalty-system-gophermart/internal/models"
	"accrual-loyalty-system-gophermart/internal/repository/postgres"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users interface {
	Create(ctx context.Context, login, hashedPassword string) error
	FindByLogin(ctx context.Context, login string) (*models.User, error)
}

type Repository struct {
	Users
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Users: postgres.NewUserRepository(db),
	}
}
