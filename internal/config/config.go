package config

import (
	"fmt"
	"time"

	"github.com/tnfy-link/backend/internal/id"
	"github.com/tnfy-link/core/config"
)

type HTTPConfig struct {
	Address     string   `envconfig:"HTTP__ADDRESS"`
	ProxyHeader string   `envconfig:"HTTP__PROXY_HEADER"`
	Proxies     []string `envconfig:"HTTP__PROXIES"`
}

type APIConfig struct {
	CORSAllowOrigins string `envconfig:"API__CORS_ALLOW_ORIGINS"`
}

type StorageConfig struct {
	URL string `envconfig:"STORAGE__URL"`
}

type LinksConfig struct {
	Hostname string        `envconfig:"LINKS__HOSTNAME"`
	TTL      time.Duration `envconfig:"LINKS__TTL"`
}

type IDConfig struct {
	Provider id.Provider `envconfig:"ID__PROVIDER"`
}

type Config struct {
	HTTP    HTTPConfig
	API     APIConfig
	Storage StorageConfig
	Links   LinksConfig
	ID      IDConfig
}

const (
	defaultLinksTTL = time.Hour * 24 * 7
)

func Default() Config {
	return Config{
		HTTP: HTTPConfig{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
		},
		API: APIConfig{
			CORSAllowOrigins: "",
		},
		Storage: StorageConfig{
			URL: "redis://localhost:6379/0",
		},
		Links: LinksConfig{
			Hostname: "http://localhost:3001",
			TTL:      defaultLinksTTL,
		},
		ID: IDConfig{
			Provider: id.ProviderRandom,
		},
	}
}

func New() (Config, error) {
	instance := Default()

	if err := config.Load(&instance); err != nil {
		return instance, fmt.Errorf("failed to load config: %w", err)
	}

	return instance, nil
}
