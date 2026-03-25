# Golang Health Checker

[![CI](https://github.com/essence-666/golang-health-checker/actions/workflows/ci.yml/badge.svg)](https://github.com/essence-666/golang-health-checker/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/essence-666/golang-health-checker)](https://go.dev/)
[![License](https://img.shields.io/github/license/essence-666/golang-health-checker)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/essence-666/golang-health-checker)](https://goreportcard.com/report/github.com/essence-666/golang-health-checker)
[![Docker](https://img.shields.io/badge/docker-ready-blue?logo=docker)](Dockerfile)

A lightweight, concurrent HTTP health checker written in Go. Periodically monitors a list of URLs and exposes health status, latency, and failure metrics via Prometheus.

## Features

- Concurrent health checks with configurable timeout and retries
- Prometheus metrics endpoint (`/metrics` on port 9090)
- YAML-based configuration
- Structured JSON logging (zerolog)
- Multi-stage Docker build
- Prometheus + Docker Compose monitoring stack

## Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `healthcker_up` | Gauge | `1` if target is healthy, `0` otherwise |
| `healthchecker_latency_seconds` | Gauge | Response latency in seconds |
| `healthchecker_total` | Counter | Total number of health checks |
| `healthchecker_total_failes` | Counter | Total number of failed checks |

## Quick Start

### Prerequisites

- Go 1.26+
- Docker & Docker Compose (optional)

### Run locally

```bash
# Clone the repository
git clone https://github.com/essence-666/golang-health-checker.git
cd golang-health-checker

# Edit the config
vim config.yaml

# Run
go run main.go
```

### Run with Docker Compose

```bash
# Start health-checker + Prometheus
docker compose up -d

# Health checker metrics: http://localhost:9090/metrics
# Prometheus UI:          http://localhost:9091
```

## Configuration

Edit `config.yaml`:

```yaml
checker:
  timeout: 3       # HTTP request timeout in seconds
  retries: 3       # Number of retry attempts on failure
  period: 15       # Check interval in seconds
  urls:
    - https://www.google.com
    - https://www.example.com
```

## Project Structure

```
.
├── main.go                          # Application entrypoint
├── config.yaml                      # Health check configuration
├── Dockerfile                       # Multi-stage Docker build
├── docker-compose.yaml              # Docker Compose with Prometheus
├── internal/
│   ├── api/api.go                   # HTTP server for metrics
│   ├── checker/checker.go           # Core health checking logic
│   ├── config/config.go             # YAML config parser
│   ├── logger/logger.go             # Zerolog logger setup
│   └── metrics/metrics.go           # Prometheus metrics definitions
├── prometheus/
│   └── prometheus.yml               # Prometheus scrape config
├── tests/
│   ├── unit/                        # Unit tests (go test -tags=unit)
│   └── integration/                 # Integration tests (go test -tags=integration)
└── .github/
    └── workflows/
        └── ci.yml                   # CI pipeline (lint, fmt, test, build)
```

## Testing

```bash
# Unit tests
go test -v -tags=unit ./tests/unit/...

# Integration tests
go test -v -tags=integration ./tests/integration/...

# With coverage
go test -tags=unit -coverprofile=coverage.out ./tests/unit/...
go tool cover -html=coverage.out
```

## CI Pipeline

The GitHub Actions workflow runs on every push/PR to `main`:

1. **Lint** - golangci-lint
2. **Format** - gofmt check
3. **Unit Tests** - with race detector and coverage
4. **Integration Tests** - end-to-end verification
5. **Build** - binary compilation
6. **Docker** - image build verification
