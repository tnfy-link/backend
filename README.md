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
  - [🚀 Features](#-features)
  - [🛠 Tech Stack](#-tech-stack)
  - [🏃‍♂️ Quick Start](#️-quick-start)
    - [Using Prebuilt Image](#using-prebuilt-image)
    - [Using Docker Compose with Local Build](#using-docker-compose-with-local-build)
    - [Manual Setup](#manual-setup)
  - [⚙️ Configuration](#️-configuration)
  - [📝 API Documentation](#-api-documentation)
    - [Shorten URL](#shorten-url)
    - [Redirect](#redirect)
  - [🔧 Development](#-development)
  - [🤝 Contributing](#-contributing)
  - [📄 License](#-license)

A high-performance URL shortener service built with Go, using modern technologies and best practices. This service provides fast and reliable URL shortening capabilities with Redis-based storage.


## 🚀 Features

- **High Performance**: Built with Go and Fiber framework for maximum speed
- **Redis Storage**: Fast and reliable link storage with configurable TTL
- **Base58 Encoding**: Human-friendly short URLs using Base58 encoding
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **Configurable**: Environment-based configuration for flexibility
- **Structured Logging**: Comprehensive logging with Zap logger
- **Dependency Injection**: Clean architecture using Uber's fx framework

## 🛠 Tech Stack

- **Language**: Go 1.23+
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber)
- **Storage**: Redis
- **Logging**: Uber Zap
- **DI Framework**: Uber fx
- **Containerization**: Docker

## 🏃‍♂️ Quick Start

### Using Prebuilt Image

The fastest way to get started is using our prebuilt image from GitHub Container Registry:

1. Create a docker-compose.yml:
    ```yaml
    services:
      server:
        image: ghcr.io/tnfy-link/backend:latest
        ports:
          - "3000:3000"
        environment:
          - HTTP__ADDRESS=:3000
          - STORAGE__URL=redis://redis:6379/0
          - LINKS__HOSTNAME=http://localhost:3000
          - LINKS__TTL=168h
        depends_on:
          - redis
      
      redis:
        image: redis:7-alpine
        volumes:
          - redis-data:/data

    volumes:
      redis-data:
    ```

2. Start the services:
    ```bash
    docker compose up -d
    ```

The service will be available at `http://localhost:3000`

### Using Docker Compose with Local Build

1. Clone the repository:
    ```bash
    git clone https://github.com/tnfy-link/server.git
    cd server
    ```

2. Start the services:
    ```bash
    docker compose up --build -d
    ```

### Manual Setup

1. Prerequisites:
   - Go 1.23 or later
   - Redis server

2. Install dependencies:
    ```bash
    go mod download
    ```

3. Configure environment variables:
    ```bash
    cp .env.example .env
    # Edit .env with your configuration
    ```

4. Run the server:
    ```bash
    go run main.go
    ```

## ⚙️ Configuration

Configuration is done through environment variables:

| Variable          | Description                       | Default                    |
| ----------------- | --------------------------------- | -------------------------- |
| `HTTP__ADDRESS`   | HTTP server listen address        | `:3000`                    |
| `STORAGE__URL`    | Redis connection URL              | `redis://localhost:6379/0` |
| `LINKS__HOSTNAME` | Base hostname for generated links | `http://localhost:3000`    |
| `LINKS__TTL`      | Time-to-live for shortened links  | `168h`                     |

## 📝 API Documentation

### Shorten URL
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
    "createdAt": "2024-12-09T12:53:12.76979501Z"
  }
}
```

### Redirect
```http
GET /{id}
```
Redirects to the original URL if found.

## 🔧 Development

The project uses a modular architecture with the following structure:

- `/internal/core` - Core components (HTTP server, Redis, Logger)
- `/internal/links` - Link shortening business logic
- `/internal/config` - Configuration management

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

---
Built with ❤️ using Go and [Codeium](https://codeium.com).
