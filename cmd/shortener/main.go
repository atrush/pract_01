package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v6"

	"github.com/atrush/pract_01.git/internal/api"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage/inmemory"
)

type Config struct {
	ServerPort string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("error getting config: %s\n", err.Error())
	}

	db := inmemory.NewStorage()
	svc, err := service.NewShortURLService(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	handler := api.NewHandler(svc, cfg.BaseURL)

	server := api.NewServer(cfg.ServerPort, *handler)
	log.Fatal(server.Run())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutdown server: %s\n", err.Error())
	}
}
