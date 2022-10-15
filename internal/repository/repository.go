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

type Repository struct {
	Users
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Users: postgres.NewUserRepository(db),
	}
}
