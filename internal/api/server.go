package api

import (
	"context"
	"errors"
	"fmt"
	mgrpc "github.com/atrush/pract_01.git/internal/grpc"
	pb "github.com/atrush/pract_01.git/internal/grpc/proto"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/pkg"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

//  Server implements http server
type Server struct {
	httpServer http.Server
	grpcServer *grpc.Server
	cfg        *pkg.Config
}

//  NewServer return new server
func NewServer(cfg *pkg.Config, db storage.Storage) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("error server initiation: config is nil")
	}

	svcSht, err := service.NewShortURLService(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации handler:%w", err)
	}
	svcUser, err := service.NewUserService(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации handler:%w", err)
	}

	handler, err := NewHandler(svcSht, svcUser, cfg.BaseURL, cfg.TrustedSubnet)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации handler:%w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterURLsServer(grpcServer, mgrpc.NewURLServer(svcSht, cfg.BaseURL))

	return &Server{
		httpServer: http.Server{
			Addr:    cfg.ServerPort,
			Handler: NewRouter(handler, cfg.Debug),
		},
		grpcServer: grpcServer,
		cfg:        cfg,
	}, nil
}

//  Run starts GRPC server
func (s *Server) RunGRPC() error {
	listen, err := net.Listen("tcp", ":3201")
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(listen)
}

//  Run starts http server
//  if config EnableHTTPS true runs in HTTPS mode
func (s *Server) RunHTTP() error {
	if s.cfg.EnableHTTPS {
		certPath, keyPath, err := pkg.GetCertX509Files()
		if err != nil {
			return fmt.Errorf("error serve ssl:%w", err)
		}
		return handleServerCloseErr(s.httpServer.ListenAndServeTLS(certPath, keyPath))
	}

	return handleServerCloseErr(s.httpServer.ListenAndServe())
}

//  returns error if error is not http.ErrServerClosed
func handleServerCloseErr(err error) error {
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server closed with: %w", err)
	}

	return nil
}

//  Shutdown sutdown http server
func (s *Server) ShutdownHTTP(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

//  Shutdown sutdown http server
func (s *Server) ShutdownGRPC() {
	s.grpcServer.GracefulStop()
}
