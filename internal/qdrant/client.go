package qdrant

import (
	"github.com/caarlos0/env/v11"
	"github.com/qdrant/go-client/qdrant"
)

type config struct {
	Host string `env:"QDRANT_HOST" envDefault:"localhost"`
	Port int    `env:"QDRANT_PORT" envDefault:"6334"`
}

func NewClient() (*qdrant.Client, error) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return qdrant.NewClient(&qdrant.Config{
		Host: cfg.Host,
		Port: cfg.Port,
	})
}
