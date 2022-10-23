package handlers

import (
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/resources"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) GetOrders(c echo.Context) error {
	user := c.Get("user").(*models.User)

	orders, err := h.services.Order.All(c.Request().Context(), user)
	if err != nil {
		log.Error().Err(err).Msg("Getting all user orders")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if len(orders) > 0 {
		return c.JSON(http.StatusOK, resources.MakeOrderResourceCollection(orders))
	} else {
		return c.NoContent(http.StatusNoContent)
	}
}
