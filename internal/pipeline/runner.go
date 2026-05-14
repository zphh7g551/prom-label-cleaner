package pipeline

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/prom-label-cleaner/internal/exporter"
	"github.com/prom-label-cleaner/internal/pruner"
	"github.com/prom-label-cleaner/internal/reporter"
)

// RunResult holds the outcome of a single pipeline execution.
type RunResult struct {
	MetricFamiliesTotal int
	LabelsRemoved       int
	DryRun              bool
}

// Run executes the full pipeline: scrape → parse → detect → prune → export → report.
func (p *Pipeline) Run(ctx context.Context) (*RunResult, error) {
	body, err := p.scraper.Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("scrape: %w", err)
	}

	families, err := p.parser.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	for _, fam := range families {
		for _, mf := range fam.GetMetric() {
			for _, lp := range mf.GetLabel() {
				p.detector.Observe(fam.GetName(), lp.GetName(), lp.GetValue())
			}
		}
	}

	pruneCfg := pruner.BuildConfigFromDetector(p.detector)
	pr := pruner.New(pruneCfg)

	result := &RunResult{
		MetricFamiliesTotal: len(families),
		DryRun:              p.cfg.DryRun,
	}

	var pruned []*dto.MetricFamily
	if !p.cfg.DryRun {
		pruned, result.LabelsRemoved = pr.Prune(families)
	} else {
		pruned = families
	}

	var out io.Writer = os.Stdout
	exp := exporter.NewWithOptions(exporter.WithOutput(out))
	if err := exp.Write(pruned); err != nil {
		return nil, fmt.Errorf("export: %w", err)
	}

	rep := reporter.New(os.Stderr)
	rep.Report(p.detector.Stats())

	return result, nil
}
