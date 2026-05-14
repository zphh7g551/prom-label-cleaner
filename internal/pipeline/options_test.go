package pipeline_test

import (
	"testing"
	"time"

	"github.com/prom-label-cleaner/internal/pipeline"
)

func TestNewWithOptionsAppliesURL(t *testing.T) {
	p, err := pipeline.NewWithOptions(
		pipeline.WithTargetURL("http://localhost:9090/metrics"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestNewWithOptionsAppliesThreshold(t *testing.T) {
	p, err := pipeline.NewWithOptions(
		pipeline.WithTargetURL("http://localhost:9090/metrics"),
		pipeline.WithCardinalityThreshold(50),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestNewWithOptionsAppliesTimeout(t *testing.T) {
	p, err := pipeline.NewWithOptions(
		pipeline.WithTargetURL("http://localhost:9090/metrics"),
		pipeline.WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestNewWithOptionsDryRun(t *testing.T) {
	p, err := pipeline.NewWithOptions(
		pipeline.WithTargetURL("http://localhost:9090/metrics"),
		pipeline.WithDryRun(true),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestNewWithOptionsBearerToken(t *testing.T) {
	p, err := pipeline.NewWithOptions(
		pipeline.WithTargetURL("http://localhost:9090/metrics"),
		pipeline.WithBearerToken("secret-token"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}
