package pipeline

import "time"

// Option is a functional option for Pipeline configuration.
type Option func(*Config)

// WithTargetURL sets the scrape target URL.
func WithTargetURL(url string) Option {
	return func(c *Config) {
		c.Scraper.TargetURL = url
	}
}

// WithBearerToken sets the bearer token used for scraping.
func WithBearerToken(token string) Option {
	return func(c *Config) {
		c.Scraper.BearerToken = token
	}
}

// WithTimeout sets the HTTP client timeout for scraping.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Scraper.Timeout = d
	}
}

// WithCardinalityThreshold sets the cardinality threshold for the detector.
func WithCardinalityThreshold(n int) Option {
	return func(c *Config) {
		c.Detector.Threshold = n
	}
}

// WithDryRun enables or disables dry-run mode.
func WithDryRun(enabled bool) Option {
	return func(c *Config) {
		c.DryRun = enabled
	}
}

// NewWithOptions creates a Pipeline applying the given functional options
// on top of DefaultConfig.
func NewWithOptions(opts ...Option) (*Pipeline, error) {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return New(cfg)
}
