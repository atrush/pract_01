package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler Handler) *chi.Mux {
	r := chi.NewRouter()
	compress := middleware.NewCompressor(5, "/*")
	r.Use(compress.Handler)
	r.Get("/{shortID}", handler.GetURLHandler)
	r.Post("/", handler.SaveURLHandler)
	r.Post("/api/shorten", handler.SaveURLJSONHandler)

	return r
}
