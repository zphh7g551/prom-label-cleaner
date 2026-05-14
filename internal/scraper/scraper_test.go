package scraper_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prom-label-cleaner/internal/scraper"
)

func TestFetchSuccess(t *testing.T) {
	expected := `# HELP go_goroutines Number of goroutines\ngo_goroutines 42\n`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expected))
	}))
	defer server.Close()

	cfg := scraper.DefaultConfig()
	cfg.TargetURL = server.URL
	s := scraper.New(cfg)

	data, err := s.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestFetchNonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	cfg := scraper.DefaultConfig()
	cfg.TargetURL = server.URL
	s := scraper.New(cfg)

	_, err := s.Fetch()
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestFetchBearerTokenHeader(t *testing.T) {
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := scraper.DefaultConfig()
	cfg.TargetURL = server.URL
	cfg.BearerToken = "secret-token"
	s := scraper.New(cfg)

	_, _ = s.Fetch()
	if receivedAuth != "Bearer secret-token" {
		t.Errorf("expected 'Bearer secret-token', got %q", receivedAuth)
	}
}

func TestFetchInvalidURL(t *testing.T) {
	cfg := scraper.DefaultConfig()
	cfg.TargetURL = "://invalid-url"
	s := scraper.New(cfg)

	_, err := s.Fetch()
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}
