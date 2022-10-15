package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "An address and port for server to start")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "An address of DB connection")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "An address of the Accrual System")

	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Parsing env")
		return nil
	}

	log.Debug().Interface("config", cfg).Send()
	return &cfg
}
