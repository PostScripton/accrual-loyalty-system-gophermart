package middlewares

import (
	"github.com/labstack/echo/v4"
)

func AcceptJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "application/json" {
			return echo.NewHTTPError(400, "Available Content-Type is application/json")
		}

		return next(c)
	}
}
