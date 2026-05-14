package exporter_test

import (
	"strings"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"google.golang.org/protobuf/proto"

	"github.com/yourorg/prom-label-cleaner/internal/exporter"
)

func parseFamily(t *testing.T, text string) []*dto.MetricFamily {
	t.Helper()
	var families []*dto.MetricFamily
	dec := expfmt.NewDecoder(strings.NewReader(text), expfmt.NewFormat(expfmt.TypeTextPlain))
	for {
		mf := &dto.MetricFamily{}
		if err := dec.Decode(mf); err != nil {
			break
		}
		families = append(families, mf)
	}
	return families
}

const sampleMetrics = `# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET"} 42
`

func TestWriteToStringTextFormat(t *testing.T) {
	families := parseFamily(t, sampleMetrics)
	if len(families) == 0 {
		t.Fatal("expected at least one family")
	}

	cfg := exporter.DefaultConfig()
	exp := exporter.New(cfg)

	out, err := exp.WriteToString(families)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "http_requests_total") {
		t.Errorf("output missing metric name, got: %s", out)
	}
}

func TestWriteSkipsNilFamily(t *testing.T) {
	name := "up"
	mfType := dto.MetricType_GAUGE
	families := []*dto.MetricFamily{
		nil,
		{
			Name: proto.String(name),
			Type: &mfType,
			Metric: []*dto.Metric{
				{Gauge: &dto.Gauge{Value: proto.Float64(1)}},
			},
		},
	}

	cfg := exporter.DefaultConfig()
	exp := exporter.New(cfg)

	out, err := exp.WriteToString(families)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "up") {
		t.Errorf("expected 'up' in output, got: %s", out)
	}
}

func TestWriteEmptyFamilies(t *testing.T) {
	cfg := exporter.DefaultConfig()
	exp := exporter.New(cfg)

	out, err := exp.WriteToString(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}
