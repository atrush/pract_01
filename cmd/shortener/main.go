package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/atrush/pract_01.git/internal/api"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/infile"
	"github.com/atrush/pract_01.git/internal/storage/inmemory"
	"github.com/atrush/pract_01.git/pkg"
)

func main() {

	cfg, err := pkg.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := getInitDB(*cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	svc, err := service.NewShortURLService(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	handler := api.NewHandler(svc, cfg.BaseURL)

	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer ctxCancel()

	server := api.NewServer(cfg.ServerPort, *handler)
	log.Fatal(server.Run())

	<-ctx.Done()
	ctxCancel()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("error shutdown server: %s\n", err.Error())
	}
}

func getInitDB(cfg pkg.Config) (storage.URLStorer, error) {
	if cfg.FileStoragePath != "" {
		db, err := infile.NewFileStorage(cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return inmemory.NewStorage(), nil
}
