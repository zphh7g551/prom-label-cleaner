package parser_test

import (
	"testing"

	"github.com/prom-label-cleaner/internal/parser"
)

const sampleMetrics = `# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",status="200"} 1027
http_requests_total{method="POST",status="500"} 3
# HELP go_goroutines Number of goroutines
# TYPE go_goroutines gauge
go_goroutines 17
`

func TestParseValidMetrics(t *testing.T) {
	families, err := parser.Parse([]byte(sampleMetrics))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(families) != 2 {
		t.Errorf("expected 2 metric families, got %d", len(families))
	}
	if _, ok := families["http_requests_total"]; !ok {
		t.Error("expected 'http_requests_total' family not found")
	}
	if _, ok := families["go_goroutines"]; !ok {
		t.Error("expected 'go_goroutines' family not found")
	}
}

func TestParseEmptyInput(t *testing.T) {
	families, err := parser.Parse([]byte{})
	if err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
	if len(families) != 0 {
		t.Errorf("expected 0 families, got %d", len(families))
	}
}

func TestLabelNamesForFamily(t *testing.T) {
	families, err := parser.Parse([]byte(sampleMetrics))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	mf, ok := families["http_requests_total"]
	if !ok {
		t.Fatal("metric family not found")
	}

	labels := parser.LabelNamesForFamily(mf)
	if _, ok := labels["method"]; !ok {
		t.Error("expected label 'method' not found")
	}
	if _, ok := labels["status"]; !ok {
		t.Error("expected label 'status' not found")
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}
