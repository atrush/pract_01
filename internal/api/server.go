package api

import (
	"context"
	"errors"
	"net/http"
)

type Server struct {
	httpServer http.Server
}

func NewServer(port string, handler Handler) *Server {

	return &Server{
		httpServer: http.Server{
			Addr:    port,
			Handler: NewRouter(handler),
		},
	}
}

func (s *Server) Run() error {

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.ListenAndServe(); err == http.ErrServerClosed {
		return errors.New("http server not runned")
	}
	return s.httpServer.Shutdown(ctx)
}
