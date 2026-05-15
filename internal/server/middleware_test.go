package server_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prom-label-cleaner/internal/server"
)

func newTestLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, nil))
}

func TestLoggingMiddlewareLogsRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := server.LoggingMiddleware(logger, handler)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	logged := buf.String()
	if !strings.Contains(logged, "request") {
		t.Errorf("expected log to contain 'request', got: %s", logged)
	}
	if !strings.Contains(logged, "/metrics") {
		t.Errorf("expected log to contain path, got: %s", logged)
	}
}

func TestLoggingMiddlewareCapturesStatusCode(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	mw := server.LoggingMiddleware(logger, handler)
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	if !strings.Contains(buf.String(), "404") {
		t.Errorf("expected log to contain status 404, got: %s", buf.String())
	}
}

func TestRecoveryMiddlewareHandlesPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	mw := server.RecoveryMiddleware(logger, handler)
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
	if !strings.Contains(buf.String(), "panic recovered") {
		t.Errorf("expected panic log, got: %s", buf.String())
	}
}

func TestRecoveryMiddlewarePassesNormalRequests(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := server.RecoveryMiddleware(logger, handler)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
