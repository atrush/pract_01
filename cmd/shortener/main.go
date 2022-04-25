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

	ctx, cancel := context.WithCancel(context.Background())
	finishedChan := make(chan struct{})

	db, err := getDB(ctx, finishedChan, *cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	server, err := api.NewServer(cfg, db)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		if err := server.RunGRPC(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	go func() {
		if err := server.RunHTTP(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	// graceful shutdown
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc
	log.Println("accepted sigint, shutting down")

	if err := server.ShutdownHTTP(ctx); err != nil {
		log.Fatalf("error shutdown server: %s\n", err.Error())
	}

	server.ShutdownGRPC()

	// waiting if db must wait ending of async tasks
	if db.WaitAsyncTasksEnded() {
		log.Println("waiting all tasks finish")
		cancel()
		<-finishedChan
		log.Println("all tasks finished")
	}
}

//  getDB returns initialized storage
//  psql storage if dsn not empty, else memory storage
func getDB(ctx context.Context, finishedChan chan struct{}, cfg pkg.Config) (storage.Storage, error) {
	log.Println("dsn: " + cfg.DatabaseDSN)
	//  postgress storage
	if cfg.DatabaseDSN != "" {
		db, err := psql.NewStorage(ctx, finishedChan, cfg.DatabaseDSN)
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
