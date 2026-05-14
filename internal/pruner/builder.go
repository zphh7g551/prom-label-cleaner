package pruner

import (
	"github.com/your-org/prom-label-cleaner/internal/cardinality"
)

// BuildConfigFromDetector creates a pruner Config by querying a Detector
// for all metric families that have high-cardinality labels and collecting
// the offending label names per family.
func BuildConfigFromDetector(d *cardinality.Detector) Config {
	cfg := DefaultConfig()
	for family, stats := range d.AllStats() {
		var highCard []string
		for label, st := range stats {
			if d.IsHighCardinality(family, label) {
				highCard = append(highCard, label)
				_ = st
			}
		}
		if len(highCard) > 0 {
			cfg.LabelsToPrune[family] = highCard
		}
	}
	return cfg
}
