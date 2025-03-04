package id

import (
	"github.com/redis/go-redis/v9"
	"github.com/tnfy-link/backend/internal/id/provider"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"id",
	fx.Provide(NewGenerator),
	fx.Provide(fx.Annotate(newProvider, fx.As(new(idProvider))), fx.Private),
)

func newProvider(config Config, storage *redis.Client) idProvider {
	switch config.Provider {
	case ProviderRandom:
		return provider.NewRandomGenerator()
	case ProviderCombined:
		return provider.NewCombinedGenerator(storage)
	default:
		panic("unknown provider")
	}
}
