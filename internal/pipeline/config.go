package pipeline

import (
	"fmt"
	"time"
)

// Config holds all configuration for the pipeline.
type Config struct {
	TargetURL            string
	BearerToken          string
	Timeout              time.Duration
	CardinalityThreshold int
	DryRun               bool
	OutputFormat         string
}

// Validate checks that the Config has all required fields and valid values.
func (c *Config) Validate() error {
	if c.TargetURL == "" {
		return fmt.Errorf("pipeline config: target URL must not be empty")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("pipeline config: timeout must be positive, got %s", c.Timeout)
	}
	if c.CardinalityThreshold <= 0 {
		return fmt.Errorf("pipeline config: cardinality threshold must be positive, got %d", c.CardinalityThreshold)
	}
	validFormats := map[string]bool{
		"text":         true,
		"openmetrics":  true,
	}
	if c.OutputFormat != "" && !validFormats[c.OutputFormat] {
		return fmt.Errorf("pipeline config: unsupported output format %q", c.OutputFormat)
	}
	return nil
}
