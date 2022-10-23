package middlewares

import (
	"errors"
	"fmt"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/services"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type Auth struct {
	Services *services.Services
}

func (middleware *Auth) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := middleware.parseAuthorizationHeader(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		token, err := middleware.parseJWT(tokenString)

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user, err := middleware.Services.User.FindByLogin(c.Request().Context(), claims["sub"].(string))
			if err != nil || user == nil {
				log.Error().Err(err).Msg("Finding user after valid JWT")
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			c.Set("user", user)
		} else {
			log.Error().Err(err).Msg("Parse JWT")
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		return next(c)
	}
}

func (middleware *Auth) parseAuthorizationHeader(c echo.Context) (string, error) {
	authorization := c.Request().Header.Get("Authorization")
	if authorization == "" {
		return "", errors.New("authorizations is required")
	}
	return strings.Split(authorization, " ")[1], nil
}

func (middleware *Auth) parseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(middleware.Services.Auth.GetSecret()), nil
	})
}
