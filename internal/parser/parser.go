package parser

import (
	"bytes"
	"fmt"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// MetricFamily is an alias for the Prometheus data model type.
type MetricFamily = dto.MetricFamily

// Parse decodes raw Prometheus text-format metrics into MetricFamily objects.
func Parse(data []byte) (map[string]*MetricFamily, error) {
	reader := bytes.NewReader(data)
	decoder := expfmt.NewDecoder(reader, expfmt.NewFormat(expfmt.TypeTextPlain))

	families := make(map[string]*MetricFamily)

	for {
		mf := &dto.MetricFamily{}
		err := decoder.Decode(mf)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("decoding metrics: %w", err)
		}
		if mf.Name != nil {
			families[*mf.Name] = mf
		}
	}

	return families, nil
}

// LabelNamesForFamily returns the set of all label names used across a MetricFamily.
func LabelNamesForFamily(mf *MetricFamily) map[string]struct{} {
	names := make(map[string]struct{})
	for _, m := range mf.Metric {
		for _, lp := range m.Label {
			if lp.Name != nil {
				names[*lp.Name] = struct{}{}
			}
		}
	}
	return names
}
