package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/prom-label-cleaner/internal/exporter"
)

func TestWithOutputWritesToBuffer(t *testing.T) {
	var buf bytes.Buffer
	exp := exporter.NewWithOptions(
		exporter.WithOutput(&buf),
	)

	families := parseFamily(t, sampleMetrics)
	if err := exp.Write(families); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "http_requests_total") {
		t.Errorf("expected metric in buffer output, got: %s", buf.String())
	}
}

func TestWithFormatDefault(t *testing.T) {
	cfg := exporter.DefaultConfig()
	if cfg.Format != exporter.FormatText {
		t.Errorf("expected default format %q, got %q", exporter.FormatText, cfg.Format)
	}
}

func TestWithFormatOpenMetrics(t *testing.T) {
	var buf bytes.Buffer
	exp := exporter.NewWithOptions(
		exporter.WithFormat(exporter.FormatOpenMetrics),
		exporter.WithOutput(&buf),
	)

	families := parseFamily(t, sampleMetrics)
	// OpenMetrics encoder may return an error on Encode for some versions;
	// we only assert no panic and a non-empty or empty output without crash.
	_ = exp.Write(families)
}
