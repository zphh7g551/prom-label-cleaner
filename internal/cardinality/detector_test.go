package cardinality

import (
	"fmt"
	"testing"
)

func TestObserveAndStats(t *testing.T) {
	d := NewDetector(DefaultDetectorConfig())

	for i := 0; i < 10; i++ {
		d.Observe("pod", fmt.Sprintf("pod-%d", i))
	}
	for i := 0; i < 3; i++ {
		d.Observe("env", fmt.Sprintf("env-%d", i))
	}

	stats := d.Stats()
	if len(stats) != 2 {
		t.Fatalf("expected 2 label stats, got %d", len(stats))
	}
	// Stats should be sorted descending by unique values.
	if stats[0].Name != "pod" {
		t.Errorf("expected first stat to be 'pod', got %q", stats[0].Name)
	}
	if stats[0].UniqueValues != 10 {
		t.Errorf("expected 10 unique values for 'pod', got %d", stats[0].UniqueValues)
	}
	if stats[1].UniqueValues != 3 {
		t.Errorf("expected 3 unique values for 'env', got %d", stats[1].UniqueValues)
	}
}

func TestHighCardinality(t *testing.T) {
	cfg := DetectorConfig{Threshold: 5}
	d := NewDetector(cfg)

	for i := 0; i < 200; i++ {
		d.Observe("request_id", fmt.Sprintf("req-%d", i))
	}
	for i := 0; i < 3; i++ {
		d.Observe("status", fmt.Sprintf("%d", i))
	}

	hc := d.HighCardinality()
	if len(hc) != 1 {
		t.Fatalf("expected 1 high-cardinality label, got %d", len(hc))
	}
	if hc[0].Name != "request_id" {
		t.Errorf("expected high-cardinality label to be 'request_id', got %q", hc[0].Name)
	}
}

func TestSummaryNoHighCardinality(t *testing.T) {
	d := NewDetector(DefaultDetectorConfig())
	d.Observe("env", "production")
	d.Observe("env", "staging")

	summary := d.Summary()
	expected := fmt.Sprintf("no high-cardinality labels detected (threshold: %d)", DefaultDetectorConfig().Threshold)
	if summary != expected {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestSummaryWithHighCardinality(t *testing.T) {
	cfg := DetectorConfig{Threshold: 2}
	d := NewDetector(cfg)
	for i := 0; i < 10; i++ {
		d.Observe("trace_id", fmt.Sprintf("t-%d", i))
	}

	summary := d.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestDuplicateValuesNotCounted(t *testing.T) {
	d := NewDetector(DefaultDetectorConfig())
	for i := 0; i < 50; i++ {
		d.Observe("region", "us-east-1")
	}

	stats := d.Stats()
	if len(stats) != 1 {
		t.Fatalf("expected 1 stat, got %d", len(stats))
	}
	if stats[0].UniqueValues != 1 {
		t.Errorf("expected 1 unique value, got %d", stats[0].UniqueValues)
	}
}
