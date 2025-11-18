# Prometheus Metrics Documentation

This document details all Prometheus metrics exposed by the application, including their types, descriptions, labels, and implementation details.

## Custom Application Metrics

### Links Module

#### links_created_total

- **Type**: Counter
- **Description**: Total number of shortened links created
- **Subsystem**: links
- **Labels**: None
- **Implementation**: [`internal/links/metrics.go:23-28`](internal/links/metrics.go:23-28)
- **Incremented by**: [`internal/links/metrics.go:41-43`](internal/links/metrics.go:41-43)

#### links_errors_total

- **Type**: CounterVec (with labels)
- **Description**: Total number of link errors
- **Subsystem**: links
- **Labels**: `reason` (possible values: `validation_failed`, `id_generation_failed`, `storage_failed`)
- **Implementation**: [`internal/links/metrics.go:30-36`](internal/links/metrics.go:30-36)
- **Incremented by**: [`internal/links/metrics.go:45-47`](internal/links/metrics.go:45-47)

### Stats Module

#### stats_redirects_total

- **Type**: Counter
- **Description**: Total number of link redirects
- **Subsystem**: stats
- **Labels**: None
- **Implementation**: [`internal/stats/metrics.go:14-20`](internal/stats/metrics.go:14-20)
- **Incremented by**: [`internal/stats/metrics.go:24-26`](internal/stats/metrics.go:24-26)

## HTTP Metrics

### http_requests_total

- **Type**: Counter
- **Description**: Count of all HTTP requests by status code, method and path
- **Labels**: `method`, `path`, `status_code`

### http_request_duration_seconds

- **Type**: Histogram
- **Description**: Duration of all HTTP requests by status code, method and path
- **Labels**: `method`, `path`, `status_code`

### http_requests_in_progress_total

- **Type**: Gauge
- **Description**: Count of all requests in progress
- **Labels**: `method`

## Standard Go Runtime Metrics

The Prometheus client library automatically exposes Go runtime metrics including:

- Go garbage collection metrics (`go_gc_*`)
- Go memory statistics (`go_memstats_*`)
- Go runtime information (`go_info`, `go_goroutines`, `go_threads`, etc.)
- Process metrics (`process_*`)

## Dependencies

- **Prometheus Client**: `github.com/prometheus/client_golang`
- **Fiberprometheus**: `github.com/ansrivas/fiberprometheus/v2`
- **Promauto**: Used for automatic registration of metrics
