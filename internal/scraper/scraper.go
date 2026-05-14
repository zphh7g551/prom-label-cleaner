package scraper

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config holds configuration for the Prometheus scraper.
type Config struct {
	TargetURL  string
	Timeout    time.Duration
	BearerToken string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout: 10 * time.Second,
	}
}

// Scraper fetches raw Prometheus metrics from a target endpoint.
type Scraper struct {
	cfg    Config
	client *http.Client
}

// New creates a new Scraper with the given configuration.
func New(cfg Config) *Scraper {
	return &Scraper{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// Fetch retrieves the raw metrics text from the configured target URL.
func (s *Scraper) Fetch() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, s.cfg.TargetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if s.cfg.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.cfg.BearerToken)
	}

	req.Header.Set("Accept", "text/plain; version=0.0.4")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return data, nil
}
