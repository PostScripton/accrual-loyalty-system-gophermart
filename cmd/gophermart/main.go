package main

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/config"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/clients"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository/postgres"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/server"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

// go run cmd/gophermart/main.go -d=postgres://homestead:secret@localhost:5432/accrual_loyalty_system -r=localhost:8081

// JWTSecret todo get jwt secret from .env
const JWTSecret = "OCf7CyrgOfnXT1udxqOOVfC5QSBnkGau"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	cfg := config.NewConfig()

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := postgres.NewPostgres(mainCtx, cfg.DatabaseURI)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	defer db.Close()

	client := clients.NewAccrualSystemClient(cfg.AccrualSystemAddress)
	repo := repository.NewRepository(db)
	newServices := services.NewServices(repo, client, JWTSecret)

	s := server.NewServer(cfg.RunAddress, newServices)

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return s.Run()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return s.Shutdown(context.Background())
	})
	g.Go(func() error {
		if err = newServices.Order.RunPollingStatuses(mainCtx); err != nil {
			log.Error().Err(err).Msg("Failed polling statuses")
			return err
		}
		return nil
	})

	if err = g.Wait(); err != nil {
		log.Info().Msg("The application is shutdown")
	}
}
