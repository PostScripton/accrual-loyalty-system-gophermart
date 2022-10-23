package handlers

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) CreateOrder(c echo.Context) error {
	numberBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("Read body on creating order")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	number := string(numberBytes)
	numberInt, err := strconv.Atoi(number)
	if err != nil {
		log.Error().Err(err).Msg("Converting body bytes to int")
		return echo.NewHTTPError(http.StatusBadRequest, "The number is required to be passed as a plain text")
	}

	if !h.services.Luhn.Valid(numberInt) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid order numberBytes")
	}

	user := c.Get("user").(*models.User)
	order, err := h.services.Order.FindByNumber(context.TODO(), number, user)
	if err != nil {
		log.Error().Err(err).Msg("Find order by numberBytes")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if order != nil {
		if order.UserID == user.ID {
			return echo.NewHTTPError(http.StatusOK, "The order has already been taken")
		} else {
			return echo.NewHTTPError(http.StatusConflict, "The order number is already registered by another user")
		}
	}

	if _, err = h.services.Order.Create(context.TODO(), number, user); err != nil {
		log.Error().Err(err).Msg("Create order")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusAccepted, "The order has been accepted")
}
