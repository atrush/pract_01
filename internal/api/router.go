package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "text/html",
		"text/css",
		"text/plain",
		"text/javascript",
		"application/javascript",
		"application/x-javascript",
		"application/json",
		"application/atom+xml",
		"application/rss+xml",
		"image/svg+xml"))
	r.Use(gzipReaderHandle)

	r.Get("/{shortID}", handler.GetURLHandler)
	r.Post("/", handler.SaveURLHandler)
	r.Post("/api/shorten", handler.SaveURLJSONHandler)

	return r
}
