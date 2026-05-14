package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/prom-label-cleaner/internal/cardinality"
)

func sampleStats() map[string]cardinality.LabelStats {
	return map[string]cardinality.LabelStats{
		"http_requests_total": {
			LabelCardinality:     map[string]int{"status": 3, "user_id": 500},
			HighCardinalityLabels: []string{"user_id"},
		},
		"db_query_duration": {
			LabelCardinality:     map[string]int{"query": 10},
			HighCardinalityLabels: []string{},
		},
	}
}

func TestReportTextOutput(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)

	if err := r.Report(sampleStats()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "http_requests_total") {
		t.Error("expected metric name in output")
	}
	if !strings.Contains(out, "user_id") {
		t.Error("expected label name in output")
	}
	if !strings.Contains(out, "YES") {
		t.Error("expected high cardinality marker in output")
	}
	if !strings.Contains(out, "METRIC") {
		t.Error("expected header in output")
	}
}

func TestReportTextEmptyStats(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)

	if err := r.Report(map[string]cardinality.LabelStats{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No metrics observed") {
		t.Error("expected empty message")
	}
}

func TestReportJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatJSON)

	if err := r.Report(sampleStats()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "http_requests_total") {
		t.Error("expected metric name in JSON output")
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Error("expected JSON to start with {")
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	r := New(nil, FormatText)
	if r.out == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
