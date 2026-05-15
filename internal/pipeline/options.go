package pipeline

import (
	"fmt"
	"time"
)

// Option is a functional option for configuring a Pipeline via New.
type Option func(*Config) error

// WithTargetURL sets the Prometheus metrics endpoint to scrape.
func WithTargetURL(url string) Option {
	return func(c *Config) error {
		if url == "" {
			return fmt.Errorf("target URL must not be empty")
		}
		c.TargetURL = url
		return nil
	}
}

// WithBearerToken sets the bearer token used for authenticated scrapes.
func WithBearerToken(token string) Option {
	return func(c *Config) error {
		c.BearerToken = token
		return nil
	}
}

// WithTimeout sets the HTTP client timeout for scrape requests.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) error {
		if d <= 0 {
			return fmt.Errorf("timeout must be positive")
		}
		c.Timeout = d
		return nil
	}
}

// WithCardinalityThreshold sets the minimum unique-value count that
// causes a label to be considered high-cardinality.
func WithCardinalityThreshold(n int) Option {
	return func(c *Config) error {
		if n < 0 {
			return fmt.Errorf("cardinality threshold must be non-negative")
		}
		c.CardinalityThreshold = n
		return nil
	}
}

// WithDryRun enables dry-run mode: detection is performed but no labels
// are removed from the exported output.
func WithDryRun(v bool) Option {
	return func(c *Config) error {
		c.DryRun = v
		return nil
	}
}

// WithScrapeInterval sets the interval used by the Scheduler.
func WithScrapeInterval(d time.Duration) Option {
	return func(c *Config) error {
		if d <= 0 {
			return fmt.Errorf("scrape interval must be positive")
		}
		c.ScrapeInterval = d
		return nil
	}
}
