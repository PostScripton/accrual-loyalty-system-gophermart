package main

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/config"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository/postgres"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/server"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// go run cmd/gophermart/main.go -d=postgres://homestead:secret@localhost:5432/accrual_loyalty_system

// JWTSecret todo get jwt secret from .env
const JWTSecret = "OCf7CyrgOfnXT1udxqOOVfC5QSBnkGau"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	cfg := config.NewConfig()

	db, err := postgres.NewPostgres(context.Background(), cfg.DatabaseURI)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	newServices := services.NewServices(repo, JWTSecret)

	s := server.NewServer(cfg.RunAddress, newServices)
	s.Run()
}
