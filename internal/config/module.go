package config

import (
	"github.com/tnfy-link/server/internal/core/http"
	"github.com/tnfy-link/server/internal/core/redis"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	fx.Provide(New),
	fx.Provide(func(c Config) http.Config {
		return http.Config{
			Address: c.Http.Address,
		}
	}),
	fx.Provide(func(c Config) redis.Config {
		return redis.Config{
			URL: c.Storage.URL,
		}
	}),
)
