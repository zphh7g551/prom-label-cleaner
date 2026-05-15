package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Addr:         ":9090",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// Server wraps an HTTP server that exposes cleaned metrics.
type Server struct {
	cfg    Config
	server *http.Server
}

// New creates a new Server with the given config and handler.
func New(cfg Config, handler http.Handler) (*Server, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("server config: %w", err)
	}
	return &Server{
		cfg: cfg,
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}, nil
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server listen: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Addr returns the configured listen address.
func (s *Server) Addr() string {
	return s.cfg.Addr
}
