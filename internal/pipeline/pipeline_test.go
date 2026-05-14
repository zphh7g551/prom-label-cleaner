package pipeline_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prom-label-cleaner/internal/cardinality"
	"github.com/prom-label-cleaner/internal/pipeline"
	"github.com/prom-label-cleaner/internal/scraper"
)

const sampleMetrics = `# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/a"} 1
http_requests_total{method="GET",path="/b"} 2
http_requests_total{method="POST",path="/c"} 3
`

func newTestServer(body string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestRunDryRun(t *testing.T) {
	srv := newTestServer(sampleMetrics, http.StatusOK)
	defer srv.Close()

	cfg := pipeline.DefaultConfig()
	cfg.Scraper = scraper.Config{TargetURL: srv.URL}
	cfg.Detector = cardinality.DetectorConfig{Threshold: 2}
	cfg.DryRun = true

	p, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	res, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if res.Cleaned != res.Raw {
		t.Error("dry-run: cleaned output should equal raw output")
	}
}

func TestRunPrunesHighCardinality(t *testing.T) {
	srv := newTestServer(sampleMetrics, http.StatusOK)
	defer srv.Close()

	cfg := pipeline.DefaultConfig()
	cfg.Scraper = scraper.Config{TargetURL: srv.URL}
	cfg.Detector = cardinality.DetectorConfig{Threshold: 2}
	cfg.DryRun = false

	p, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	res, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(res.Pruned) == 0 {
		t.Error("expected at least one pruned label")
	}

	for _, label := range res.Pruned {
		if strings.Contains(res.Cleaned, label) {
			t.Errorf("pruned label %q still present in cleaned output", label)
		}
	}
}

func TestRunFetchError(t *testing.T) {
	cfg := pipeline.DefaultConfig()
	cfg.Scraper = scraper.Config{TargetURL: "http://127.0.0.1:0"}

	p, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = p.Run(context.Background())
	if err == nil {
		t.Fatal("expected error from unreachable target, got nil")
	}
}
