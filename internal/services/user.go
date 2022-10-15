package services

import (
	"context"
	"errors"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrLoginTaken = errors.New("this login is already taken")

type UserService struct {
	repo repository.Users
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) Create(ctx context.Context, login, password string) (*models.User, error) {
	user, err := us.repo.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, ErrLoginTaken
	}

	hashedPassword, err := us.hashPassword(password)
	if err != nil {
		return nil, err
	}

	if err = us.repo.Create(ctx, login, hashedPassword); err != nil {
		return nil, err
	}

	return us.repo.FindByLogin(ctx, login)
}

func (us *UserService) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	return us.repo.FindByLogin(ctx, login)
}

func (us *UserService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
