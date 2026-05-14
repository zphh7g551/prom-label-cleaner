package pipeline

import (
	"context"
	"fmt"

	"github.com/prom-label-cleaner/internal/cardinality"
	"github.com/prom-label-cleaner/internal/parser"
	"github.com/prom-label-cleaner/internal/pruner"
	"github.com/prom-label-cleaner/internal/scraper"
)

// Config holds configuration for the pipeline.
type Config struct {
	Scraper   scraper.Config
	Detector  cardinality.DetectorConfig
	DryRun    bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Scraper:  scraper.DefaultConfig(),
		Detector: cardinality.DefaultDetectorConfig(),
		DryRun:   false,
	}
}

// Result holds the output of a single pipeline run.
type Result struct {
	Raw     string
	Cleaned string
	Pruned  []string
}

// Pipeline orchestrates scraping, detection, and pruning.
type Pipeline struct {
	cfg     Config
	scraper *scraper.Scraper
}

// New creates a new Pipeline.
func New(cfg Config) (*Pipeline, error) {
	s, err := scraper.New(cfg.Scraper)
	if err != nil {
		return nil, fmt.Errorf("pipeline: init scraper: %w", err)
	}
	return &Pipeline{cfg: cfg, scraper: s}, nil
}

// Run executes the full scrape → detect → prune pipeline.
func (p *Pipeline) Run(ctx context.Context) (*Result, error) {
	raw, err := p.scraper.Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("pipeline: fetch: %w", err)
	}

	families, err := parser.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("pipeline: parse: %w", err)
	}

	det := cardinality.NewDetector(p.cfg.Detector)
	for _, fam := range families {
		for _, lbl := range parser.LabelNamesForFamily(fam) {
			for _, m := range fam.GetMetric() {
				for _, lp := range m.GetLabel() {
					if lp.GetName() == lbl {
						det.Observe(fam.GetName(), lbl, lp.GetValue())
					}
				}
			}
		}
	}

	pruneCfg := pruner.BuildConfigFromDetector(det)
	var prunedLabels []string
	for metric, labels := range pruneCfg.Rules {
		for _, l := range labels {
			prunedLabels = append(prunedLabels, metric+"."+l)
		}
	}

	if p.cfg.DryRun {
		return &Result{Raw: raw, Cleaned: raw, Pruned: prunedLabels}, nil
	}

	pr := pruner.New(pruneCfg)
	cleaned, err := pr.Prune(families)
	if err != nil {
		return nil, fmt.Errorf("pipeline: prune: %w", err)
	}

	return &Result{Raw: raw, Cleaned: cleaned, Pruned: prunedLabels}, nil
}
