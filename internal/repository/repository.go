package repository

import (
	"accrual-loyalty-system-gophermart/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users interface {
}

type Repository struct {
	Users
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Users: postgres.NewUserRepository(db),
	}
}
