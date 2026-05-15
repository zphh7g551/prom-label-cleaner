package server

import (
	"testing"
	"time"
)

func TestValidateOK(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}

func TestValidateMissingAddr(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Addr = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty addr")
	}
}

func TestValidateZeroReadTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ReadTimeout = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero read timeout")
	}
}

func TestValidateZeroWriteTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.WriteTimeout = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero write timeout")
	}
}

func TestValidateZeroIdleTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.IdleTimeout = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero idle timeout")
	}
}

func TestDefaultConfigValues(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Addr != ":9090" {
		t.Errorf("expected default addr :9090, got %q", cfg.Addr)
	}
	if cfg.ReadTimeout != 10*time.Second {
		t.Errorf("expected 10s read timeout, got %v", cfg.ReadTimeout)
	}
}
