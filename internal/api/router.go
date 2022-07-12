package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//  NewRouter init server routes.
func NewRouter(handler *Handler, debug bool) *chi.Mux {
	r := chi.NewRouter()

	//  compress middlewares
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

	//  pprof routes
	if debug {
		r.Mount("/debug", middleware.Profiler())
	}

	//  route for allowed subnets
	r.Group(func(r chi.Router) {
		r.Use(handler.subnet.Middleware)
		r.Get("/api/internal/stats", handler.Stats)
	})

	//  json auth routes
	r.Group(func(r chi.Router) {
		r.Use(handler.auth.Middleware)
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/shorten/batch", handler.SaveBatch)
		r.Post("/api/shorten", handler.SaveURLJSONHandler)
		r.Delete("/api/user/urls", handler.DeleteBatch)
	})

	// auth routes
	r.Group(func(r chi.Router) {
		r.Use(handler.auth.Middleware)
		r.Get("/ping", handler.Ping)
		r.Get("/api/user/urls", handler.GetUserUrls)
		r.Get("/{shortID}", handler.GetURLHandler)
		r.Post("/", handler.SaveURLHandler)
	})

	return r
}
