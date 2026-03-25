package api

import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func ExportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9090", nil)
}