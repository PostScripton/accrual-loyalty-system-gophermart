package handlers

import (
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/services"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}
