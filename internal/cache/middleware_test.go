package cache_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prom-label-cleaner/internal/cache"
)

func TestMiddlewareCachesResponse(t *testing.T) {
	hits := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("metrics_data"))
	})

	c := cache.New()
	mw := cache.Middleware(c, handler)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	rr1 := httptest.NewRecorder()
	mw.ServeHTTP(rr1, req)

	rr2 := httptest.NewRecorder()
	mw.ServeHTTP(rr2, req)

	if hits != 1 {
		t.Errorf("expected handler called once, got %d", hits)
	}
	if rr2.Body.String() != "metrics_data" {
		t.Errorf("expected cached body, got %q", rr2.Body.String())
	}
}

func TestMiddlewareExpiredEntryRefetches(t *testing.T) {
	hits := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fresh_data"))
	})

	c := cache.NewWithOptions(cache.WithTTL(50 * time.Millisecond))
	mw := cache.Middleware(c, handler)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	rr1 := httptest.NewRecorder()
	mw.ServeHTTP(rr1, req)

	time.Sleep(100 * time.Millisecond)

	rr2 := httptest.NewRecorder()
	mw.ServeHTTP(rr2, req)

	if hits != 2 {
		t.Errorf("expected handler called twice after expiry, got %d", hits)
	}
}

func TestMiddlewareNonGetPassthrough(t *testing.T) {
	hits := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
	})

	c := cache.New()
	mw := cache.Middleware(c, handler)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/metrics", nil)
		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)
	}

	if hits != 3 {
		t.Errorf("expected POST requests to bypass cache, got %d hits", hits)
	}
}
