package server

import (
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/middlewares"
)

func (s *Server) registerRoutes() {
	apiGroup := s.core.Group("/api", middlewares.AcceptJSON)
	userGroup := apiGroup.Group("/user")

	authMiddleware := &middlewares.Auth{Services: s.services}
	authGroup := userGroup.Group("", authMiddleware.Handle)
	simpleAuthGroup := s.core.Group("/api/user", authMiddleware.Handle)

	userGroup.POST("/register", s.handler.Register)
	userGroup.POST("/login", s.handler.Login)

	simpleAuthGroup.POST("/orders", s.handler.CreateOrder)
	simpleAuthGroup.GET("/orders", s.handler.GetOrders)

	simpleAuthGroup.GET("/balance", s.handler.Balance)
	authGroup.POST("/balance/withdraw", s.handler.Withdraw)
	simpleAuthGroup.GET("/withdrawals", s.handler.GetWithdrawals)
}
