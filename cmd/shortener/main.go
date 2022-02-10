package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/atrush/pract_01.git/internal/api"
	"github.com/atrush/pract_01.git/internal/storage/infile"
	"github.com/atrush/pract_01.git/internal/storage/psql"
	"github.com/atrush/pract_01.git/pkg"
)

func main() {

	cfg, err := pkg.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := infile.NewFileStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	psDB, err := getInitPsDB(*cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer psDB.Close()

	server, err := api.NewServer(cfg, db, psDB)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(server.Run())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutdown server: %s\n", err.Error())
	}
}

func getInitPsDB(cfg pkg.Config) (*psql.Storage, error) {
	if cfg.DatabaseDSN != "" {
		db, err := psql.NewStorage(cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	return nil, nil
}
