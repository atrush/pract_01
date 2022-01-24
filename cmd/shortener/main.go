package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/atrush/pract_01.git/internal/api"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/infile"
	"github.com/atrush/pract_01.git/internal/storage/inmemory"
)

func main() {
	var cfg api.Config
	api.ReadEnvConfig(&cfg)

	db, err := getInitDB(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

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

func getInitDB(cfg api.Config) (storage.URLStorer, error) {
	if cfg.FileStoragePath != "" {
		db, err := infile.NewFileStorage(cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return inmemory.NewStorage(), nil
}
