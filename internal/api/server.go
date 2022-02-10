package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/psql"
	"github.com/atrush/pract_01.git/pkg"
)

type Server struct {
	httpServer http.Server
}

// Return new server
func NewServer(cfg *pkg.Config, db storage.Storage, psDB *psql.Storage) (*Server, error) {
	svc, err := service.NewService(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации handler:%w", err)
	}

	handler, err := NewHandler(svc, psDB, cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации handler:%w", err)
	}

	return &Server{
		httpServer: http.Server{
			Addr:    cfg.ServerPort,
			Handler: NewRouter(handler),
		},
	}, nil
}

// Start server
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown server
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.ListenAndServe(); err == http.ErrServerClosed {
		return errors.New("http server not runned")
	}

	return s.httpServer.Shutdown(ctx)
}
