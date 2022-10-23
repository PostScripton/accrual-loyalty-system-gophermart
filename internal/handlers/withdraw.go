package handlers

import (
	"fmt"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type withdrawal struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required"`
}

func (h *Handler) Withdraw(c echo.Context) error {
	user := c.Get("user").(*models.User)

	w := new(withdrawal)
	if err := c.Bind(w); err != nil {
		log.Error().Err(err).Msg("Bind withdrawalModel")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := validator.New().Struct(w); err != nil {
		err := err.(validator.ValidationErrors)[0]
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("The %s is %s", err.Field(), err.Tag()))
	}

	number, err := strconv.Atoi(w.Order)
	if err != nil {
		log.Error().Err(err).Msg("Number string to int conversion")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !h.services.Luhn.Valid(number) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	if user.Balance < w.Sum {
		return echo.NewHTTPError(http.StatusPaymentRequired, "Insufficient funds")
	}

	withdrawalModel := &models.Withdrawal{
		UserID: user.ID,
		Order:  w.Order,
		Sum:    w.Sum,
	}
	if err = h.services.Withdrawal.Create(c.Request().Context(), withdrawalModel, user); err != nil {
		log.Error().Err(err).Msg("Create withdrawalModel")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
