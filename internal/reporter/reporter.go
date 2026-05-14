package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/prom-label-cleaner/internal/cardinality"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Reporter writes cardinality summaries to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to the given writer with the specified format.
// If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Report writes a summary of cardinality stats for all tracked metrics.
func (r *Reporter) Report(stats map[string]cardinality.LabelStats) error {
	switch r.format {
	case FormatJSON:
		return r.reportJSON(stats)
	default:
		return r.reportText(stats)
	}
}

func (r *Reporter) reportText(stats map[string]cardinality.LabelStats) error {
	if len(stats) == 0 {
		_, err := fmt.Fprintln(r.out, "No metrics observed.")
		return err
	}

	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC\tLABEL\tCARDINALITY\tHIGH")

	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, metric := range keys {
		ls := stats[metric]
		labels := make([]string, 0, len(ls.LabelCardinality))
		for l := range ls.LabelCardinality {
			labels = append(labels, l)
		}
		sort.Strings(labels)
		for _, label := range labels {
			card := ls.LabelCardinality[label]
			high := ""
			for _, h := range ls.HighCardinalityLabels {
				if h == label {
					high = "YES"
					break
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", metric, label, card, high)
		}
	}
	return w.Flush()
}

func (r *Reporter) reportJSON(stats map[string]cardinality.LabelStats) error {
	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(r.out, "{")
	for i, metric := range keys {
		ls := stats[metric]
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(r.out, "  %q: {\"high_cardinality_labels\": %q, \"label_count\": %d}%s\n",
			metric, ls.HighCardinalityLabels, len(ls.LabelCardinality), comma)
	}
	fmt.Fprintln(r.out, "}")
	return nil
}
