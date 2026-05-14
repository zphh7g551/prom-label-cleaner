// Package exporter writes pruned Prometheus metrics to an output destination.
package exporter

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/prometheus/common/expfmt"
	dto "github.com/prometheus/client_model/go"
)

// Format represents the output format for exported metrics.
type Format string

const (
	FormatText       Format = "text"
	FormatOpenMetrics Format = "openmetrics"
)

// Config holds configuration for the Exporter.
type Config struct {
	Format Format
	Output io.Writer
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Format: FormatText,
		Output: os.Stdout,
	}
}

// Exporter serialises a slice of MetricFamily values to the configured output.
type Exporter struct {
	cfg Config
}

// New creates a new Exporter using the provided Config.
func New(cfg Config) *Exporter {
	return &Exporter{cfg: cfg}
}

// Write encodes families and writes them to the configured output.
func (e *Exporter) Write(families []*dto.MetricFamily) error {
	var enc expfmt.Encoder

	switch e.cfg.Format {
	case FormatOpenMetrics:
		enc = expfmt.NewEncoder(e.cfg.Output, expfmt.NewOpenMetricsFormat(expfmt.OpenMetricsType))
	default:
		enc = expfmt.NewEncoder(e.cfg.Output, expfmt.NewFormat(expfmt.TypeTextPlain))
	}

	for _, mf := range families {
		if mf == nil {
			continue
		}
		if err := enc.Encode(mf); err != nil {
			return fmt.Errorf("exporter: encode %q: %w", mf.GetName(), err)
		}
	}
	return nil
}

// WriteToString encodes families and returns the result as a string.
func (e *Exporter) WriteToString(families []*dto.MetricFamily) (string, error) {
	var buf bytes.Buffer
	orig := e.cfg.Output
	e.cfg.Output = &buf
	defer func() { e.cfg.Output = orig }()

	if err := e.Write(families); err != nil {
		return "", err
	}
	return buf.String(), nil
}
