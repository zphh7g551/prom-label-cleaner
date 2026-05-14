package cardinality

import (
	"fmt"
	"sort"
)

// LabelStats holds cardinality information for a single label name.
type LabelStats struct {
	Name        string
	UniqueValues int
	SampleValues []string
}

// DetectorConfig configures the cardinality detector.
type DetectorConfig struct {
	// Threshold is the maximum number of unique values a label may have
	// before it is considered high-cardinality.
	Threshold int
}

// DefaultDetectorConfig returns a DetectorConfig with sensible defaults.
func DefaultDetectorConfig() DetectorConfig {
	return DetectorConfig{
		Threshold: 100,
	}
}

// Detector analyses label value cardinality in a set of metric samples.
type Detector struct {
	cfg    DetectorConfig
	labels map[string]map[string]struct{}
}

// NewDetector creates a new Detector with the given configuration.
func NewDetector(cfg DetectorConfig) *Detector {
	return &Detector{
		cfg:    cfg,
		labels: make(map[string]map[string]struct{}),
	}
}

// Observe records a label name/value pair for later analysis.
func (d *Detector) Observe(name, value string) {
	if _, ok := d.labels[name]; !ok {
		d.labels[name] = make(map[string]struct{})
	}
	d.labels[name][value] = struct{}{}
}

// Stats returns cardinality statistics for every observed label, sorted by
// unique value count descending.
func (d *Detector) Stats() []LabelStats {
	stats := make([]LabelStats, 0, len(d.labels))
	for name, values := range d.labels {
		samples := make([]string, 0, min(5, len(values)))
		for v := range values {
			samples = append(samples, v)
			if len(samples) == 5 {
				break
			}
		}
		stats = append(stats, LabelStats{
			Name:        name,
			UniqueValues: len(values),
			SampleValues: samples,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].UniqueValues > stats[j].UniqueValues
	})
	return stats
}

// HighCardinality returns only the labels whose unique value count exceeds the
// configured threshold.
func (d *Detector) HighCardinality() []LabelStats {
	var result []LabelStats
	for _, s := range d.Stats() {
		if s.UniqueValues > d.cfg.Threshold {
			result = append(result, s)
		}
	}
	return result
}

// Summary returns a human-readable summary string.
func (d *Detector) Summary() string {
	hc := d.HighCardinality()
	if len(hc) == 0 {
		return fmt.Sprintf("no high-cardinality labels detected (threshold: %d)", d.cfg.Threshold)
	}
	return fmt.Sprintf("%d high-cardinality label(s) detected (threshold: %d)", len(hc), d.cfg.Threshold)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
