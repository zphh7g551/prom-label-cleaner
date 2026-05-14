package pruner

import (
	"testing"

	"github.com/your-org/prom-label-cleaner/internal/cardinality"
)

func TestBuildConfigFromDetector(t *testing.T) {
	cfg := cardinality.DefaultDetectorConfig()
	cfg.Threshold = 3
	d := cardinality.NewDetector(cfg)

	// Simulate observations: "user_id" will exceed threshold, "method" will not.
	for _, v := range []string{"u1", "u2", "u3", "u4"} {
		d.Observe("http_requests_total", "user_id", v)
	}
	for _, v := range []string{"GET", "POST"} {
		d.Observe("http_requests_total", "method", v)
	}

	pruneCfg := BuildConfigFromDetector(d)

	labels, ok := pruneCfg.LabelsToPrune["http_requests_total"]
	if !ok {
		t.Fatal("expected http_requests_total to be in prune config")
	}

	found := false
	for _, l := range labels {
		if l == "user_id" {
			found = true
		}
		if l == "method" {
			t.Error("method should not be marked high-cardinality")
		}
	}
	if !found {
		t.Error("expected user_id to be in prune list")
	}
}

func TestBuildConfigFromDetectorNoHighCardinality(t *testing.T) {
	cfg := cardinality.DefaultDetectorConfig()
	cfg.Threshold = 100
	d := cardinality.NewDetector(cfg)

	d.Observe("go_goroutines", "instance", "localhost:9090")

	pruneCfg := BuildConfigFromDetector(d)
	if len(pruneCfg.LabelsToPrune) != 0 {
		t.Errorf("expected empty prune config, got %v", pruneCfg.LabelsToPrune)
	}
}
