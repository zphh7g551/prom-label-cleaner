package cardinality

import "sync"

// DetectorConfig holds configuration for the Detector.
type DetectorConfig struct {
	// Threshold is the number of unique label values above which a label is
	// considered high-cardinality.
	Threshold int
}

// DefaultDetectorConfig returns a DetectorConfig with sensible defaults.
func DefaultDetectorConfig() DetectorConfig {
	return DetectorConfig{Threshold: 100}
}

// LabelStats holds cardinality statistics for a single label.
type LabelStats struct {
	UniqueValues int
}

// Detector tracks label value cardinality across metric families.
type Detector struct {
	mu     sync.RWMutex
	cfg    DetectorConfig
	// stats: family -> label -> set of unique values
	stats  map[string]map[string]map[string]struct{}
}

// NewDetector creates a new Detector with the given config.
func NewDetector(cfg DetectorConfig) *Detector {
	return &Detector{
		cfg:   cfg,
		stats: make(map[string]map[string]map[string]struct{}),
	}
}

// Observe records a label value observation for the given family and label.
func (d *Detector) Observe(family, label, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.stats[family]; !ok {
		d.stats[family] = make(map[string]map[string]struct{})
	}
	if _, ok := d.stats[family][label]; !ok {
		d.stats[family][label] = make(map[string]struct{})
	}
	d.stats[family][label][value] = struct{}{}
}

// Stats returns the LabelStats for a specific family and label.
func (d *Detector) Stats(family, label string) (LabelStats, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if lblMap, ok := d.stats[family]; ok {
		if vals, ok := lblMap[label]; ok {
			return LabelStats{UniqueValues: len(vals)}, true
		}
	}
	return LabelStats{}, false
}

// IsHighCardinality reports whether the given label in the given family
// exceeds the configured threshold.
func (d *Detector) IsHighCardinality(family, label string) bool {
	st, ok := d.Stats(family, label)
	if !ok {
		return false
	}
	return st.UniqueValues > d.cfg.Threshold
}

// AllStats returns a snapshot of all tracked statistics.
// The returned map is family -> label -> LabelStats.
func (d *Detector) AllStats() map[string]map[string]LabelStats {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make(map[string]map[string]LabelStats, len(d.stats))
	for family, lblMap := range d.stats {
		out[family] = make(map[string]LabelStats, len(lblMap))
		for label, vals := range lblMap {
			out[family][label] = LabelStats{UniqueValues: len(vals)}
		}
	}
	return out
}

// Summary returns a map of family -> label -> LabelStats for labels that
// exceed the cardinality threshold.
func (d *Detector) Summary() map[string]map[string]LabelStats {
	all := d.AllStats()
	out := make(map[string]map[string]LabelStats)
	for family, lblMap := range all {
		for label, st := range lblMap {
			if st.UniqueValues > d.cfg.Threshold {
				if out[family] == nil {
					out[family] = make(map[string]LabelStats)
				}
				out[family][label] = st
			}
		}
	}
	return out
}
