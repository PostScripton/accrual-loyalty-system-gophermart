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

type credentials struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Register(c echo.Context) error {
	creds := new(credentials)
	if err := c.Bind(creds); err != nil {
		return err
	}

	if err := validator.New().Struct(creds); err != nil {
		err := err.(validator.ValidationErrors)[0]
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("The %s is %s", err.Field(), err.Tag()))
	}

	user, err := h.services.User.Create(c.Request().Context(), creds.Login, creds.Password)
	if err != nil {
		if errors.Is(err, services.ErrLoginTaken) {
			return echo.NewHTTPError(http.StatusConflict, "This login is already taken")
		}

		log.Error().Err(err).Msg("Create user")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	token, err := h.services.Auth.LoginByUser(user)
	if err != nil {
		log.Error().Err(err).Msg("Login user")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return c.NoContent(http.StatusOK)
}
