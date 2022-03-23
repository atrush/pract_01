package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/pkg"
)

//  Server implements http server
type Server struct {
	httpServer http.Server
}

//  NewServer return new server
func NewServer(cfg *pkg.Config, db storage.Storage) (*Server, error) {
	svcSht, err := service.NewShortURLService(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации handler:%w", err)
	}
	svcUser, err := service.NewUserService(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации handler:%w", err)
	}

	handler, err := NewHandler(svcSht, svcUser, cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации handler:%w", err)
	}

	return &Server{
		httpServer: http.Server{
			Addr:    cfg.ServerPort,
			Handler: NewRouter(handler, cfg.Debug),
		},
	}, nil
}

//  Run starts http server
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

//  Shutdown sutdown http server
func (s *Server) Shutdown(ctx context.Context) error {
	//  check server not off
	if err := s.httpServer.ListenAndServe(); err == http.ErrServerClosed {
		return errors.New("http server not runned")
	}

	return s.httpServer.Shutdown(ctx)
}
