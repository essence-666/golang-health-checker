package main

import (
	"time"

	"github.com/essence-666/golang-health-checker/internal/api"
	"github.com/essence-666/golang-health-checker/internal/checker"
	"github.com/essence-666/golang-health-checker/internal/config"
	"github.com/essence-666/golang-health-checker/internal/logger"
	"github.com/essence-666/golang-health-checker/internal/metrics"
)

func main() {

	logger := logger.GetLogger()

	logger.Info().Msg("Parse the configuration yaml file")
	configChecker, err := config.ReadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load config, terminating...")
	}

	m := metrics.NewMetric()

	c := checker.NewChecker(logger, configChecker.Timeout, configChecker.Retries, m)

	go api.ExportMetrics()

	ticker := time.NewTicker(time.Duration(configChecker.Period) * time.Second)
	for range ticker.C {
		c.RunChecker(configChecker.URLs)
	}
}
