package pruner

import (
	"fmt"
	"strings"

	dto "github.com/prometheus/client_model/go"
)

// Config holds configuration for the pruner.
type Config struct {
	// LabelsToPrune is a map of metric family name to label names to remove.
	LabelsToPrune map[string][]string
}

// DefaultConfig returns a Config with empty pruning rules.
func DefaultConfig() Config {
	return Config{
		LabelsToPrune: make(map[string][]string),
	}
}

// Pruner removes high-cardinality labels from metric families.
type Pruner struct {
	cfg Config
}

// New creates a new Pruner with the given config.
func New(cfg Config) *Pruner {
	return &Pruner{cfg: cfg}
}

// Prune removes configured labels from the given metric families.
// It returns a new slice of metric families with labels removed.
func (p *Pruner) Prune(families []*dto.MetricFamily) ([]*dto.MetricFamily, error) {
	result := make([]*dto.MetricFamily, 0, len(families))
	for _, mf := range families {
		if mf == nil {
			continue
		}
		labels, ok := p.cfg.LabelsToPrune[mf.GetName()]
		if !ok || len(labels) == 0 {
			result = append(result, mf)
			continue
		}
		pruned, err := pruneFamily(mf, labels)
		if err != nil {
			return nil, fmt.Errorf("pruning family %q: %w", mf.GetName(), err)
		}
		result = append(result, pruned)
	}
	return result, nil
}

// pruneFamily returns a copy of mf with the given label names removed from all metrics.
func pruneFamily(mf *dto.MetricFamily, labelNames []string) (*dto.MetricFamily, error) {
	remove := make(map[string]struct{}, len(labelNames))
	for _, l := range labelNames {
		remove[strings.TrimSpace(l)] = struct{}{}
	}

	name := mf.GetName()
	help := mf.GetHelp()
	typ := mf.GetType()

	newMetrics := make([]*dto.Metric, 0, len(mf.GetMetric()))
	for _, m := range mf.GetMetric() {
		filtered := make([]*dto.LabelPair, 0, len(m.GetLabel()))
		for _, lp := range m.GetLabel() {
			if _, skip := remove[lp.GetName()]; !skip {
				filtered = append(filtered, lp)
			}
		}
		newM := *m
		newM.Label = filtered
		newMetrics = append(newMetrics, &newM)
	}

	return &dto.MetricFamily{
		Name:   &name,
		Help:   &help,
		Type:   &typ,
		Metric: newMetrics,
	}, nil
}
