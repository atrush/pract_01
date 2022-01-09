package main

import (
	"log"
	"net/http"

	"github.com/atrush/pract_01.git/internal/handlers"
	"github.com/atrush/pract_01.git/internal/storage"
)

func main() {
	handler := handlers.Handler{DB: storage.NewStorage()}
	http.HandleFunc("/", handler.RequestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
