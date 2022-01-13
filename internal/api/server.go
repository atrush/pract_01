package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer http.Server
}

func NewServer(port string, handler Handler) *Server {
	r := chi.NewRouter()
	r.Get("/*", handler.GetURLHandler)
	r.Post("/", handler.SaveURLHandler)

	return &Server{
		httpServer: http.Server{
			Addr:    port,
			Handler: r,
		},
	}
}

func (s *Server) Run() error {

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {

	return s.httpServer.Shutdown(ctx)
}
