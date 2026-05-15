package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewValidConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Addr = "127.0.0.1:0"

	s, err := New(cfg, http.DefaultServeMux)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.Addr() != cfg.Addr {
		t.Errorf("expected addr %q, got %q", cfg.Addr, s.Addr())
	}
}

func TestNewInvalidConfigReturnsError(t *testing.T) {
	cfg := Config{} // zero value — invalid
	_, err := New(cfg, http.DefaultServeMux)
	if err == nil {
		t.Fatal("expected error for invalid config, got nil")
	}
}

func TestShutdownStopsServer(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Addr = "127.0.0.1:19191"

	s, err := New(cfg, http.HandlerFunc(HealthHandler()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	started := make(chan struct{})
	go func() {
		close(started)
		_ = s.Start()
	}()

	<-started
	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		t.Errorf("shutdown error: %v", err)
	}
}

func TestHealthHandlerReturnsOK(t *testing.T) {
	rr := &fakeResponseWriter{header: make(http.Header)}
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	HealthHandler()(rr, req)
	if rr.status != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.status)
	}
	if string(rr.body) != "ok" {
		t.Errorf("expected body 'ok', got %q", string(rr.body))
	}
}

type fakeResponseWriter struct {
	header http.Header
	status int
	body   []byte
}

func (f *fakeResponseWriter) Header() http.Header        { return f.header }
func (f *fakeResponseWriter) WriteHeader(code int)       { f.status = code }
func (f *fakeResponseWriter) Write(b []byte) (int, error) {
	f.body = append(f.body, b...)
	return len(b), nil
}
var _ = fmt.Sprintf // suppress unused import
