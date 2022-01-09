package main

import (
	"log"
	"net/http"

	"github.com/atrush/pract_01.git/internal/app"
	"github.com/atrush/pract_01.git/internal/storage/mapstore"
)

func main() {
	handler := app.Handler{DB: mapstore.NewStorage()}
	http.HandleFunc("/", handler.RequestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
