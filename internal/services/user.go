package services

import (
	"accrual-loyalty-system-gophermart/internal/repository"
)

type UserService struct {
	repo repository.Users
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{
		repo: repo,
	}
}
