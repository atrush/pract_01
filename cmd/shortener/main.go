package main

import (
	"log"
	"net/http"

	"github.com/atrush/pract_01.git/internal/handlers"
	"github.com/atrush/pract_01.git/internal/storage/mapstorage"
	"github.com/go-chi/chi/v5"
)

func main() {
	handler := handlers.Handler{DB: mapstorage.NewStorage()}
	r := chi.NewRouter()
	r.Get("/{shortURL}", handler.GetURLHandler)
	r.Post("/", handler.SaveURLHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
