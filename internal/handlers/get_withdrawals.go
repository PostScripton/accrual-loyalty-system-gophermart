package handlers

import (
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/resources"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) GetWithdrawals(c echo.Context) error {
	user := c.Get("user").(*models.User)

	withdrawals, err := h.services.Withdrawal.All(c.Request().Context(), user)
	if err != nil {
		log.Error().Err(err).Msg("Getting all user withdrawals")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if len(withdrawals) > 0 {
		return c.JSON(http.StatusOK, resources.MakeWithdrawalResourceCollection(withdrawals))
	} else {
		return c.NoContent(http.StatusNoContent)
	}
}
