package server

import (
	"errors"
	"time"
)

// Config holds configuration for the HTTP server.
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Validate checks that the Config is valid.
func (c Config) Validate() error {
	if c.Addr == "" {
		return errors.New("addr must not be empty")
	}
	if c.ReadTimeout <= 0 {
		return errors.New("read_timeout must be positive")
	}
	if c.WriteTimeout <= 0 {
		return errors.New("write_timeout must be positive")
	}
	if c.IdleTimeout <= 0 {
		return errors.New("idle_timeout must be positive")
	}
	return nil
}
