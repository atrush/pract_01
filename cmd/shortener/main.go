package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/atrush/pract_01.git/internal/api"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage/inmemory"
)

func main() {
	db := inmemory.NewStorage()
	svc, err := service.NewShortener(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	handler := api.NewHandler(svc)

	server := api.NewServer(":8080", *handler)
	log.Fatal(server.Run())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutdown server: %s\n", err.Error())
	}
}
