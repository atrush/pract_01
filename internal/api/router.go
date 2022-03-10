package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *Handler) *chi.Mux {
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
	r.Use(handler.auth.Middleware)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/shorten/batch", handler.SaveBatch)
		r.Post("/api/shorten", handler.SaveURLJSONHandler)
		r.Delete("/api/user/urls", handler.DeleteBatch)
	})

	r.Get("/ping", handler.Ping)
	r.Get("/api/user/urls", handler.GetUserUrls)
	r.Get("/{shortID}", handler.GetURLHandler)
	r.Post("/", handler.SaveURLHandler)

	return r
}
