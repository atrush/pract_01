package api

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{shortID}", handler.GetURLHandler)
	r.Post("/", handler.SaveURLHandler)
	r.Post("/api/shorten", handler.SaveURLJSONHandler)
	return r
}
