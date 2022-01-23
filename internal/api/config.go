package api

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerPort string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func ReadEnvConfig(cfg *Config) {
	_ = env.Parse(cfg)
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080/"
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = ":8080"
	}
}
