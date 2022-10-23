package handlers

import (
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type userBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (h *Handler) Balance(c echo.Context) error {
	user := c.Get("user").(*models.User)

	withdrawn, err := h.services.Withdrawal.Sum(c.Request().Context(), user)
	if err != nil {
		log.Error().Err(err).Msg("Withdrawal sum")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, userBalance{
		Current:   user.Balance,
		Withdrawn: withdrawn,
	})
}
