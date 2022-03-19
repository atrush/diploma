package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/atrush/diploma.git/pkg"
	"github.com/atrush/diploma.git/services/auth"
	"github.com/atrush/diploma.git/storage"
	"net/http"
)

type Server struct {
	httpServer http.Server
}

func NewServer(cfg *pkg.Config, s storage.Storage) (*Server, error) {

	jwtAuth, err := auth.NewAuth(s)
	if err != nil {
		return nil, fmt.Errorf("ошибка запуска server:%w", err)
	}

	handler, err := NewHandler(jwtAuth)
	if err != nil {
		return nil, fmt.Errorf("ошибка запуска server:%w", err)
	}

	return &Server{
		httpServer: http.Server{
			Addr:    cfg.ServerAddress,
			Handler: NewRouter(handler),
		},
	}, nil
}

// Run Start server
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
