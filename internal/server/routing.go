package server

import (
	"accrual-loyalty-system-gophermart/internal/middlewares"
)

func (s *Server) registerRoutes() {
	apiGroup := s.core.Group("/api", middlewares.AcceptJSON)
	userGroup := apiGroup.Group("/user")

	userGroup.POST("/register", s.handler.Register)
	userGroup.POST("/login", s.handler.Login)
}
