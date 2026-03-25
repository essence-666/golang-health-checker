//go:build unit

package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/essence-666/golang-health-checker/internal/checker"
	"github.com/essence-666/golang-health-checker/internal/logger"
	"github.com/essence-666/golang-health-checker/internal/metrics"
	"github.com/rs/zerolog"
)

var (
	sharedLog     zerolog.Logger
	sharedMetrics *metrics.Metrics
)

func init() {
	sharedLog = logger.GetLogger()
	sharedMetrics = metrics.NewMetric()
}

func TestRunChecker_HealthyURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := checker.NewChecker(sharedLog, 5, 3, sharedMetrics)
	c.RunChecker([]string{server.URL})

	val := getGaugeValue(t, sharedMetrics.HealthStatus, server.URL)
	if val != 1 {
		t.Errorf("expected health status 1 for healthy URL, got %f", val)
	}
}

func TestRunChecker_UnhealthyURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := checker.NewChecker(sharedLog, 5, 0, sharedMetrics)
	c.RunChecker([]string{server.URL})

	val := getGaugeValue(t, sharedMetrics.HealthStatus, server.URL)
	if val != 0 {
		t.Errorf("expected health status 0 for unhealthy URL, got %f", val)
	}
}

func TestRunChecker_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := checker.NewChecker(sharedLog, 1, 0, sharedMetrics)
	c.RunChecker([]string{server.URL})

	val := getGaugeValue(t, sharedMetrics.HealthStatus, server.URL)
	if val != 0 {
		t.Errorf("expected health status 0 for timed-out URL, got %f", val)
	}
}

func TestRunChecker_MultipleURLs(t *testing.T) {
	healthy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer healthy.Close()

	unhealthy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer unhealthy.Close()

	c := checker.NewChecker(sharedLog, 5, 0, sharedMetrics)
	c.RunChecker([]string{healthy.URL, unhealthy.URL})

	healthyVal := getGaugeValue(t, sharedMetrics.HealthStatus, healthy.URL)
	if healthyVal != 1 {
		t.Errorf("expected health status 1, got %f", healthyVal)
	}

	unhealthyVal := getGaugeValue(t, sharedMetrics.HealthStatus, unhealthy.URL)
	if unhealthyVal != 0 {
		t.Errorf("expected health status 0, got %f", unhealthyVal)
	}
}
