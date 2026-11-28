# Changelog

Notable changes to the tnfy.link backend, newest first. Version headings correspond to the project's git tags; `[Unreleased]` covers changes made since the latest tag.

## [Unreleased]

### New Features

#### Observability

- **Grafana dashboard** — added a prebuilt Prometheus-backed dashboard at `deployments/grafana/dashboard.json` covering request rate, latency, link creation, errors, and redirects.
- **Alertmanager rules** — added alert rules at `deployments/alertmanager/alerts.rules.yml` for high link-creation error rate, p99 latency above one second, high HTTP request rate, and high redirect rate.

## [0.7.7] - 2026-11-19

### Maintenance

- **Go 1.26 and dependency upgrades** — upgraded the Go toolchain to 1.26, refreshed dependencies, and resolved lint findings.

## [0.7.6] - 2026-07-09

### Maintenance

- Upgraded `github.com/gofiber/fiber/v2` to 2.52.14 and `golang.org/x/net` to 0.55.0.

## [0.7.5] - 2026-05-07

### Maintenance

- Upgraded `github.com/gofiber/fiber/v2` to 2.52.13.

## [0.7.4] - 2026-04-25

### Maintenance

- Upgraded `go.opentelemetry.io/otel` to 1.41.0.

## [0.7.3] - 2026-02-25

### Maintenance

- Upgraded `github.com/gofiber/fiber/v2` to 2.52.12.

## [0.7.2] - 2026-02-10

### Maintenance

- Upgraded `github.com/gofiber/fiber/v2` to 2.52.11.

## [0.7.1] - 2025-11-25

### Maintenance

- **Go version for container images** — upgraded Go in Docker builds.
- Documented the exposed Prometheus metrics in `docs/metrics.md`.
- Upgraded `golang.org/x/crypto` to 0.45.0.

## [0.7.0] - 2025-11-12

### New Features

#### Metrics

- **Prometheus endpoint** — added `GET /metrics` exposing HTTP request count, request duration, and in-progress request metrics.
- **Application metrics** — added `links_created_total`, `links_errors_total` (with a `reason` label), and `stats_redirects_total`.

#### Deployment

- **Deploy webhook** — the release workflow can trigger a deployment webhook after publishing.

### Maintenance

- Upgraded Fiber, `golang.org/x/net`, and go-redis; refreshed GitHub Actions workflows.

## [0.6.0] - 2025-03-05

### New Features

- **Sequential-random ID provider** — added the `combined` provider, selected with `ID__PROVIDER=combined`, which mixes a Redis-backed sequential counter with random bits.

## [0.5.0] - 2025-01-17

### New Features

- **Hosted API documentation** — the generated Swagger UI is published to https://tnfy-link.github.io/backend/.
- **More reliable deploy webhook** — improved error handling when triggering deployments.

### Maintenance

- Upgraded `golang.org/x/net`.

## [0.4.0] - 2024-12-20

### New Features

- **Queued statistics processing** — redirect statistics moved into a dedicated package with queue-backed processing.
- **Optional CORS** — CORS is enabled only when `API__CORS_ALLOW_ORIGINS` is set.

### Bug Fixes

- Fixed validation error processing on API requests.

### Maintenance

- Moved shared packages into a separate module and extracted the API into its own module.
- Added major-version Docker image tags.

## [0.3.0] - 2024-12-17

### New Features

- **ID generation module** — introduced a dedicated module for short-link ID generation.
- **Link lookup endpoint** — added `GET /api/v1/links/{id}`.
- **CORS headers** — added CORS support to the API.
- **Portainer deployment** — added a Portainer-based deployment option.
- The UI switched to the API domain.

### Maintenance

- Added a Makefile and lint/test CI workflow.

## [0.2.0] - 2024-12-14

### New Features

- **Link statistics** — added redirect statistics and a `validUntil` field on links.
- **Request validation** — added validation for incoming API requests.

### Maintenance

- Switched link storage to hashed keys instead of separate keys.
- Upgraded `golang.org/x/crypto`.

## [0.1.0] - 2024-12-11

### New Features

- **Simple UI** — added a basic frontend.
- **Proxy support** — added trusted-proxy handling for HTTP requests.
- Added a README and an air configuration for local development.

## [0.0.1] - 2024-12-11

### New Features

- **Initial release** — URL shortening and redirects behind a `/v1` API prefix, a Docker image, and an Apache 2.0 license.
