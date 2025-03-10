package config

import (
	"github.com/tnfy-link/backend/internal/api"
	"github.com/tnfy-link/backend/internal/id"
	"github.com/tnfy-link/backend/internal/links"
	"github.com/tnfy-link/core/http"
	"github.com/tnfy-link/core/redis"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	fx.Provide(New),
	fx.Provide(func(c Config) http.Config {
		return http.Config{
			Address:     c.Http.Address,
			ProxyHeader: c.Http.ProxyHeader,
			Proxies:     c.Http.Proxies,
		}
	}),
	fx.Provide(func(c Config) api.Config {
		return api.Config{
			CORSAllowOrigins: c.API.CORSAllowOrigins,
		}
	}),
	fx.Provide(func(c Config) redis.Config {
		return redis.Config{
			URL: c.Storage.URL,
		}
	}),
	fx.Provide(func(c Config) links.Config {
		return links.Config{
			Hostname: c.Links.Hostname,
			TTL:      c.Links.TTL,
		}
	}),
	fx.Provide(func(c Config) id.Config {
		return id.Config{
			Provider: c.ID.Provider,
		}
	}),
)
