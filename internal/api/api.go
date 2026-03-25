package api

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ExportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("metrics server failed: %v", err)
	}
}
