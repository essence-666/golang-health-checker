//go:build integration

package integration

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/essence-666/golang-health-checker/internal/checker"
	"github.com/essence-666/golang-health-checker/internal/logger"
	"github.com/essence-666/golang-health-checker/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestHealthCheckAndMetricsExport(t *testing.T) {
	// Start a mock target server
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	log := logger.GetLogger()
	m := metrics.NewMetric()
	c := checker.NewChecker(log, 5, 1, m)

	// Run checker against mock target
	c.RunChecker([]string{target.URL})

	// Use promhttp handler directly to read metrics
	metricsServer := httptest.NewServer(promhttp.Handler())
	defer metricsServer.Close()

	resp, err := http.Get(metricsServer.URL)
	if err != nil {
		t.Fatalf("failed to fetch metrics: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	bodyStr := string(body)

	expectedMetrics := []string{
		"healthchecker_up",
		"healthchecker_latency_seconds",
		"healthchecker_checks_total",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(bodyStr, metric) {
			t.Errorf("expected metric %q not found in output", metric)
		}
	}
}
