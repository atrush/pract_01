package main

import (
	"log"

	"github.com/atrush/pract_01.git/internal/api"
	"github.com/atrush/pract_01.git/internal/storage/mapstorage"
)

func main() {
	db := mapstorage.NewStorage()
	handler := api.NewHandler(db)

	server := api.NewServer(":8080", *handler)
	log.Fatal(server.Run())
}
