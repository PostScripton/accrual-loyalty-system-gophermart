package services

import (
	"accrual-loyalty-system-gophermart/internal/models"
	"accrual-loyalty-system-gophermart/internal/repository"
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrCredentials = errors.New("credentials don't match")

type AuthService struct {
	userRepo repository.Users
	secret   string
}

func NewAuthService(repo repository.Users, secret string) *AuthService {
	return &AuthService{
		userRepo: repo,
		secret:   secret,
	}
}

func (as *AuthService) GetSecret() string {
	return as.secret
}

func (as *AuthService) LoginByUser(user *models.User) (string, error) {
	return as.generateJWT(user)
}

func (as *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := as.userRepo.FindByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrCredentials
	}

	if err = as.checkPassword(user.Password, password); err != nil {
		return "", ErrCredentials
	}

	return as.generateJWT(user)
}

func (as *AuthService) generateJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = user.Login

	tokenString, err := token.SignedString([]byte(as.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (as *AuthService) checkPassword(hashedPassword, providedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword)); err != nil {
		return err
	}

	return nil
}
