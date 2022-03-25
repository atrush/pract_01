package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/atrush/pract_01.git/internal/api"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/infile"
	"github.com/atrush/pract_01.git/internal/storage/psql"
	"github.com/atrush/pract_01.git/pkg"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version:%v\n", buildVersion)
	fmt.Printf("Build date:%v\n", buildDate)
	fmt.Printf("Build commit:%v\n", buildCommit)

	cfg, err := pkg.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := getDB(*cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	server, err := api.NewServer(cfg, db)
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

//  getDB returns initialized storage
//  psql storage if dsn not empty, else memory storage
func getDB(cfg pkg.Config) (storage.Storage, error) {
	//  postgress storage
	if cfg.DatabaseDSN != "" {
		db, err := psql.NewStorage(cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	//  memory with file storage
	db, err := infile.NewFileStorage(cfg.FileStoragePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
