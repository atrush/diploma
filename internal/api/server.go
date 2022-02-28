package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Server struct {
	httpServer http.Server
}

func NewServer() (*Server, error) {
	h, err := NewHandler()
	if err != nil {
		return nil, fmt.Errorf("error init server:%w", err)
	}
	return &Server{
		httpServer: http.Server{
			Addr:    "8081",
			Handler: NewRouter(h),
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
