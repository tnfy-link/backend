package config

import (
	"time"

	"github.com/tnfy-link/server/internal/core/config"
)

type HttpConfig struct {
	Address     string   `envconfig:"HTTP__ADDRESS"`
	ProxyHeader string   `envconfig:"HTTP__PROXY_HEADER"`
	Proxies     []string `envconfig:"HTTP__PROXIES"`
}

type StorageConfig struct {
	URL string `envconfig:"STORAGE__URL"`
}

type LinksConfig struct {
	Hostname string        `envconfig:"LINKS__HOSTNAME"`
	TTL      time.Duration `envconfig:"LINKS__TTL"`
}

type Config struct {
	Http    HttpConfig
	Storage StorageConfig
	Links   LinksConfig
}

var instance = Config{
	Http: HttpConfig{Address: "127.0.0.1:3000"},
	Storage: StorageConfig{
		URL: "redis://localhost:6379/0",
	},
	Links: LinksConfig{
		Hostname: "http://localhost:3000",
		TTL:      time.Hour * 24 * 7,
	},
}

func New() (Config, error) {
	return instance, config.Load(&instance)
}
