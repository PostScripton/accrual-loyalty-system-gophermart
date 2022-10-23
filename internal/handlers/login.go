package handlers

import (
	"errors"
	"fmt"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) Login(c echo.Context) error {
	creds := new(credentials)
	if err := c.Bind(creds); err != nil {
		return err
	}

	if err := validator.New().Struct(creds); err != nil {
		err := err.(validator.ValidationErrors)[0]
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("The %s is %s", err.Field(), err.Tag()))
	}

	token, err := h.services.Auth.Login(c.Request().Context(), creds.Login, creds.Password)
	if err != nil {
		if errors.Is(err, services.ErrCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "These credentials don't match our records")
		}

		log.Error().Err(err).Msg("Login user")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return c.NoContent(http.StatusOK)
}
