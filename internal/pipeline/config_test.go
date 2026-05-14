package pipeline

import (
	"testing"
	"time"
)

func baseConfig() *Config {
	return &Config{
		TargetURL:            "http://localhost:9090/metrics",
		Timeout:              10 * time.Second,
		CardinalityThreshold: 100,
		OutputFormat:         "text",
	}
}

func TestValidateOK(t *testing.T) {
	if err := baseConfig().Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateMissingURL(t *testing.T) {
	cfg := baseConfig()
	cfg.TargetURL = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestValidateZeroTimeout(t *testing.T) {
	cfg := baseConfig()
	cfg.Timeout = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestValidateNegativeThreshold(t *testing.T) {
	cfg := baseConfig()
	cfg.CardinalityThreshold = -1
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative threshold")
	}
}

func TestValidateUnsupportedFormat(t *testing.T) {
	cfg := baseConfig()
	cfg.OutputFormat = "yaml"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestValidateEmptyFormatIsOK(t *testing.T) {
	cfg := baseConfig()
	cfg.OutputFormat = ""
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for empty format, got: %v", err)
	}
}

func TestValidateOpenMetricsFormat(t *testing.T) {
	cfg := baseConfig()
	cfg.OutputFormat = "openmetrics"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for openmetrics format, got: %v", err)
	}
}
