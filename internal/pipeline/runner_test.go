package pipeline_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prom-label-cleaner/internal/pipeline"
)

const sampleMetrics = `# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/a"} 1
http_requests_total{method="GET",path="/b"} 2
http_requests_total{method="GET",path="/c"} 3
`

func newRunnerTestServer(body string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestRunReturnsResult(t *testing.T) {
	srv := newRunnerTestServer(sampleMetrics, http.StatusOK)
	defer srv.Close()

	p, err := pipeline.New(
		pipeline.WithTargetURL(srv.URL),
		pipeline.WithCardinalityThreshold(2),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	res, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if res.MetricFamiliesTotal == 0 {
		t.Error("expected at least one metric family")
	}
}

func TestRunDryRunNoLabelsRemoved(t *testing.T) {
	srv := newRunnerTestServer(sampleMetrics, http.StatusOK)
	defer srv.Close()

	p, err := pipeline.New(
		pipeline.WithTargetURL(srv.URL),
		pipeline.WithCardinalityThreshold(2),
		pipeline.WithDryRun(true),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	res, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if res.LabelsRemoved != 0 {
		t.Errorf("dry-run should not remove labels, got %d", res.LabelsRemoved)
	}
}

func TestRunFetchErrorPropagates(t *testing.T) {
	p, err := pipeline.New(
		pipeline.WithTargetURL("http://127.0.0.1:0"),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, runErr := p.Run(context.Background())
	if runErr == nil {
		t.Fatal("expected error from unreachable target")
	}
	if !strings.Contains(runErr.Error(), "scrape") {
		t.Errorf("expected 'scrape' in error message, got: %v", runErr)
	}
}
