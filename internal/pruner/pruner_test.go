package pruner

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

func strPtr(s string) *string { return &s }

func makeFamily(name string, labelPairs ...map[string]string) *dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, len(labelPairs))
	for _, lm := range labelPairs {
		var pairs []*dto.LabelPair
		for k, v := range lm {
			k, v := k, v
			pairs = append(pairs, &dto.LabelPair{Name: &k, Value: &v})
		}
		metrics = append(metrics, &dto.Metric{Label: pairs})
	}
	return &dto.MetricFamily{
		Name:   strPtr(name),
		Help:   strPtr("help"),
		Type:   dto.MetricType_GAUGE.Enum(),
		Metric: metrics,
	}
}

func TestPruneRemovesConfiguredLabels(t *testing.T) {
	cfg := Config{
		LabelsToPrune: map[string][]string{
			"http_requests_total": {"user_id", "session_id"},
		},
	}
	p := New(cfg)
	families := []*dto.MetricFamily{
		makeFamily("http_requests_total",
			map[string]string{"method": "GET", "user_id": "abc123", "session_id": "xyz"},
		),
	}
	result, err := p.Prune(families)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 family, got %d", len(result))
	}
	for _, m := range result[0].GetMetric() {
		for _, lp := range m.GetLabel() {
			if lp.GetName() == "user_id" || lp.GetName() == "session_id" {
				t.Errorf("label %q should have been pruned", lp.GetName())
			}
		}
	}
}

func TestPruneKeepsUnconfiguredFamilies(t *testing.T) {
	cfg := DefaultConfig()
	p := New(cfg)
	families := []*dto.MetricFamily{
		makeFamily("go_goroutines", map[string]string{"foo": "bar"}),
	}
	result, err := p.Prune(families)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !proto.Equal(families[0], result[0]) {
		t.Error("expected family to be unchanged")
	}
}

func TestPruneNilFamilySkipped(t *testing.T) {
	p := New(DefaultConfig())
	result, err := p.Prune([]*dto.MetricFamily{nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 families, got %d", len(result))
	}
}
