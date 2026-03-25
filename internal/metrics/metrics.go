package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	HealthStatus  *prometheus.GaugeVec
	Latency       *prometheus.GaugeVec
	TotalChecks   *prometheus.CounterVec
	TotalFailures *prometheus.CounterVec
}

func NewMetric() *Metrics {
	return &Metrics{
		HealthStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "healthchecker_up",
			Help: "1 if target is healthy, 0 if not",
		}, []string{"url"}),
		Latency: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "healthchecker_latency_seconds",
			Help: "request latency in seconds",
		}, []string{"url"}),
		TotalChecks: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "healthchecker_checks_total",
			Help: "total number of health checks performed",
		}, []string{"url"}),
		TotalFailures: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "healthchecker_failures_total",
			Help: "total number of failed health checks",
		}, []string{"url"}),
	}
}
