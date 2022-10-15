package services

import (
	"accrual-loyalty-system-gophermart/internal/repository"
)

type User interface {
}

type Services struct {
	User
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{
		User: NewUserService(repo.Users),
	}
}
