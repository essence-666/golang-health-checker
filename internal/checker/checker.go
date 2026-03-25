package checker

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/essence-666/golang-health-checker/internal/metrics"
	"github.com/rs/zerolog"
)

type CheckResult struct {
	URL       string
	IsHealthy bool
	Latency   time.Duration
	Error     error
}

type Checker struct {
	log          zerolog.Logger
	client       *http.Client
	retryCounter uint
	timeout      uint
	metrics      *metrics.Metrics
}

func NewChecker(log zerolog.Logger, timeout uint, retries uint, m *metrics.Metrics) *Checker {
	return &Checker{
		log:          log,
		client:       &http.Client{},
		retryCounter: retries + 1,
		timeout:      timeout,
		metrics:      m,
	}
}

func (c *Checker) checkURLHealthy(ctx context.Context, url string) (bool, time.Duration, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, 0, err
	}

	start := time.Now()

	resp, err := c.client.Do(req)
	if err != nil {
		return false, time.Since(start), err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, time.Since(start), nil
}

func (c *Checker) RunChecker(urls []string) {
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Go(func() {
			for retry := uint(1); retry <= c.retryCounter; retry++ {
				ctx, cancel := context.WithTimeout(
					context.Background(),
					time.Duration(c.timeout)*time.Second,
				)

				isHealthy, latency, err := c.checkURLHealthy(ctx, url)
				cancel()

				if ctx.Err() == context.DeadlineExceeded && retry != c.retryCounter {
					c.log.Error().Err(err).Msgf(
						"Deadline exceeded while connect to %s, retrying %d...", url, retry,
					)
					continue
				}

				if err != nil && retry != c.retryCounter {
					c.log.Error().Err(err).Msgf("Failed to check %s, retrying %d...", url, retry)
					continue
				}

				c.metrics.TotalChecks.WithLabelValues(url).Inc()

				if err != nil {
					c.metrics.TotalFailures.WithLabelValues(url).Inc()
				}

				healthStatus := 0
				if isHealthy {
					healthStatus = 1
				}

				c.metrics.HealthStatus.WithLabelValues(url).Set(float64(healthStatus))
				c.metrics.Latency.WithLabelValues(url).Set(latency.Seconds())

				c.log.Info().
					Str("url", url).
					Bool("healthy", isHealthy).
					Str("latency", latency.String()).
					Msg("Result")

				break
			}
		})
	}

	wg.Wait()
}
