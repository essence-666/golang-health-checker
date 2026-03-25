//go:build unit

package unit

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func getGaugeValue(t *testing.T, gauge *prometheus.GaugeVec, label string) float64 {
	t.Helper()

	var m dto.Metric
	if err := gauge.WithLabelValues(label).Write(&m); err != nil {
		t.Fatalf("failed to get gauge value: %v", err)
	}
	return m.GetGauge().GetValue()
}
