# tnfy.link

[![Go Report Card](https://goreportcard.com/badge/github.com/tnfy-link/backend)](https://goreportcard.com/report/github.com/tnfy-link/backend)
[![GitHub license](https://img.shields.io/github/license/tnfy-link/backend)](https://github.com/tnfy-link/backend/blob/master/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/tnfy-link/backend)](https://github.com/tnfy-link/backend)
[![Docker Image](https://img.shields.io/badge/docker-ghcr.io-blue)](https://github.com/tnfy-link/backend/pkgs/container/backend)
[![Website](https://img.shields.io/website?url=https%3A%2F%2Ftnfy.link)](https://tnfy.link)

<p align="center">
  <img src="assets/logo.png" alt="tnfy.link logo" width="200">
</p>

- [tnfy.link](#tnfylink)
  - [About The Project](#about-the-project)
  - [Features](#features)
  - [Tech Stack](#tech-stack)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Using the Pre-built Docker Image](#using-the-pre-built-docker-image)
    - [Using a Pre-built Binary](#using-a-pre-built-binary)
    - [Building from Source](#building-from-source)
  - [Usage](#usage)
    - [API Documentation](#api-documentation)
      - [Shorten URL](#shorten-url)
      - [Get Link By ID](#get-link-by-id)
      - [Get Statistics](#get-statistics)
    - [Observability](#observability)
  - [Configuration](#configuration)
  - [Contributing](#contributing)
  - [License](#license)
  - [Contact](#contact)
  - [Acknowledgments](#acknowledgments)

## About The Project

tnfy.link is a URL shortener service. This repository holds its Go backend, which creates short links, resolves them, and records redirect statistics including UTM labels. It exposes a versioned JSON API and Prometheus metrics.

The service is intended for operators running their own link shortening service and for client applications consuming the API. It is configured entirely through environment variables and ships with a Docker Compose setup for local use.

## Features

- **URL shortening** — create short links through `POST /api/v1/links`, with a configurable time-to-live
- **Link lookup** — retrieve link metadata through `GET /api/v1/links/{id}`
- **Redirect statistics** — query UTM-labelled redirect stats through `GET /api/v1/links/{id}/stats`
- **Base58 IDs** — human-friendly short IDs generated with Base58 encoding
- **Pluggable ID providers** — choose `random` or a Redis-backed `combined` sequential/random provider
- **Redis storage** — fast link storage via go-redis
- **Optional CORS** — enabled only when `API__CORS_ALLOW_ORIGINS` is set
- **Prometheus metrics** — HTTP and application metrics exposed at `/metrics`
- **Structured logging** — JSON logging with Zap, plus development logging via `DEBUG`
- **Dependency injection** — component wiring with Uber fx
- **Docker support** — local builds and deployment with Docker Compose

## Tech Stack

- **Language**: Go 1.26+
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber)
- **Storage**: [Redis](https://redis.io/) via [go-redis](https://github.com/redis/go-redis)
- **Logging**: [Uber Zap](https://github.com/uber-go/zap)
- **Dependency Injection**: [Uber fx](https://github.com/uber-go/fx)
- **Metrics**: [Prometheus client](https://github.com/prometheus/client_golang) with [fiberprometheus](https://github.com/ansrivas/fiberprometheus)
- **API Docs**: [swag](https://github.com/swaggo/swag) / Swagger UI
- **Containerization**: Docker

## Getting Started

### Prerequisites

- Docker and Docker Compose (pre-built image only)
- Redis server (pre-built binary or source build)
- Go 1.26 or later (source build only)

### Using the Pre-built Docker Image

Published images are available at `ghcr.io/tnfy-link/backend` with a `latest` tag. Start Redis and the backend:

```bash
docker network create tnfy
docker run -d --name tnfy-redis --network tnfy redis:7
docker run -d --name tnfy-backend --network tnfy -p 3000:3000 \
  -e HTTP__ADDRESS=:3000 \
  -e STORAGE__URL=redis://tnfy-redis:6379/0 \
  ghcr.io/tnfy-link/backend:latest
```

The repository's `compose.yaml` also references this image and can be used with `docker compose up -d`.

### Using a Pre-built Binary

Each tagged release publishes archives for Linux, macOS, and Windows on the [Releases page](https://github.com/tnfy-link/backend/releases). Download the archive for your platform, extract the `backend` binary, and run it:

```bash
# Example for Linux x86_64
curl -sSLO https://github.com/tnfy-link/backend/releases/latest/download/backend_Linux_x86_64.tar.gz
tar -xzf backend_Linux_x86_64.tar.gz
chmod +x backend
./backend
```

A Redis server must be reachable at the configured `STORAGE__URL` (default `redis://localhost:6379/0`).

### Building from Source

Requires Go 1.26 or later and a reachable Redis server:

```bash
git clone https://github.com/tnfy-link/backend.git
cd backend
go mod download
cp .env.example .env
# Edit .env with your configuration
go run main.go
```

To build a container image locally instead, run `docker compose up --build -d` from the cloned repository.

## Usage

### API Documentation

The API is documented with Swagger UI at the `/api/v1/docs` endpoint of the backend and at https://tnfy-link.github.io/backend/.

#### Shorten URL

```http
POST /api/v1/links
Content-Type: application/json

{
  "link": {
    "targetUrl": "https://docs.sms-gate.app"
  }
}
```

Response:
```json
{
  "link": {
    "id": "3uqH4m",
    "targetUrl": "https://docs.sms-gate.app",
    "url": "https://tnfy.link/3uqH4m",
    "createdAt": "2024-12-09T12:53:12.76979501Z",
    "validUntil": "2024-12-16T12:53:12.76979501Z"
  }
}
```

#### Get Link By ID

```http
GET /api/v1/links/{id}
```

Response:
```json
{
  "link": {
    "id": "3uqH4m",
    "targetUrl": "https://docs.sms-gate.app",
    "url": "https://tnfy.link/3uqH4m",
    "createdAt": "2024-12-09T12:53:12.76979501Z",
    "validUntil": "2024-12-16T12:53:12.76979501Z"
  }
}
```

#### Get Statistics

```http
GET /api/v1/links/{id}/stats
```

Response:
```json
{
  "stats": {
    "labels": {
      "source": {
        "google": 1
      },
      "medium": {
        "cpc": 1
      },
      "campaign": {
        "new_year": 1
      }
    },
    "total": 1
  }
}
```

### Observability

Prometheus metrics are served at `GET /metrics`. Exposed application metrics include `links_created_total`, `links_errors_total` (with a `reason` label), and `stats_redirects_total`, alongside HTTP request count, duration, and in-progress metrics. See [docs/metrics.md](docs/metrics.md) for the full list.

Ready-to-use monitoring assets are included:

- `deployments/grafana/dashboard.json` — Grafana dashboard
- `deployments/alertmanager/alerts.rules.yml` — Alertmanager rules for error rate, latency, request rate, and redirect rate

## Configuration

Configuration is read from environment variables, or from a `.env` file in the working directory. See [.env.example](.env.example) for a documented template.

| Variable                  | Description                                                 | Default                    |
| ------------------------- | ----------------------------------------------------------- | -------------------------- |
| `HTTP__ADDRESS`           | HTTP server listen address                                  | `127.0.0.1:3000`           |
| `HTTP__PROXY_HEADER`      | Trusted proxy header name                                   | `X-Forwarded-For`          |
| `HTTP__PROXIES`           | Comma-separated list of trusted proxies                     | *empty*                    |
| `API__CORS_ALLOW_ORIGINS` | CORS allowed origins; CORS is disabled when empty           | *empty*                    |
| `STORAGE__URL`            | Redis connection URL                                        | `redis://localhost:6379/0` |
| `LINKS__HOSTNAME`         | Base hostname for generated links                           | `http://localhost:3001`    |
| `LINKS__TTL`              | Time-to-live for shortened links (Go duration)              | `168h`                     |
| `ID__PROVIDER`            | ID provider (`random` or `combined`)                        | `random`                   |
| `DEBUG`                   | Enables development logging when set to any non-empty value | *unset*                    |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Run `make all` to format, lint, test, and benchmark the project before submitting.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Contact

- Support: support@tnfy.link
- Repository: https://github.com/tnfy-link/backend
- API documentation: https://tnfy-link.github.io/backend/

## Acknowledgments

- [Fiber](https://github.com/gofiber/fiber), [Uber fx](https://github.com/uber-go/fx), and [Uber Zap](https://github.com/uber-go/zap)
- [Redis](https://redis.io/) and [go-redis](https://github.com/redis/go-redis)
- [Prometheus](https://prometheus.io/) and [Grafana](https://grafana.com/)
