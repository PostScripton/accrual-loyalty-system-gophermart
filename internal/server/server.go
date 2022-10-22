package server

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/handlers"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Server struct {
	core     *echo.Echo
	handler  *handlers.Handler
	services *services.Services
	address  string
}

func NewServer(address string, services *services.Services) *Server {
	s := &Server{
		core:     echo.New(),
		handler:  handlers.NewHandler(services),
		services: services,
		address:  address,
	}
	s.registerRoutes()

	return s
}

func (s *Server) Run() error {
	log.Info().Str("address", s.address).Msg("The server has just started")
	return s.core.Start(s.address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.core.Shutdown(ctx)
}
