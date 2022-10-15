package server

import (
	"accrual-loyalty-system-gophermart/internal/middlewares"
	"accrual-loyalty-system-gophermart/internal/models"
	"fmt"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerRoutes() {
	apiGroup := s.core.Group("/api", middlewares.AcceptJSON)
	userGroup := apiGroup.Group("/user")

	userGroup.POST("/register", s.handler.Register)
	userGroup.POST("/login", s.handler.Login)

	authGroup := userGroup.Group("", (&middlewares.Auth{Services: s.services}).Handle)
	// todo this route is for testing auth middleware, remove it later
	authGroup.POST("/protected", func(c echo.Context) error {
		return c.String(200, fmt.Sprintf("protected route, user: %v", c.Get("user").(*models.User)))
	})
}
