package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{shortID}", handler.GetURLHandler)
	r.Post("/", handler.SaveURLHandler)
	r.Post("/api/shorten", handler.SaveURLJSONHandler)

	compressor := middleware.NewCompressor(5, "/*")
	r.Use(compressor.Handler)

	return r
}
