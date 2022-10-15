package server

import (
	"accrual-loyalty-system-gophermart/internal/handlers"
	"accrual-loyalty-system-gophermart/internal/services"
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

func (s *Server) Run() {
	log.Info().Str("address", s.address).Msg("The server has just started")

	if err := s.core.Start(s.address); err != nil {
		log.Fatal().Err(err).Msg("Server error occurred")
	}
}
