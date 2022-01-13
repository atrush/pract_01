package main

import (
	"log"

	"github.com/atrush/pract_01.git/internal/api"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/inmemory"
)

func main() {
	reposytory := storage.NewRepository(inmemory.NewStorage())
	svc, err := service.NewShotener(reposytory)
	if err != nil {
		log.Fatal(err.Error())
	}
	handler := api.NewHandler(svc)

	server := api.NewServer(":8080", *handler)
	log.Fatal(server.Run())
}
