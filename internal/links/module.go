package links

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"links",
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("links")
		}),
		fx.Provide(newRepository, fx.Private),
		fx.Provide(NewService),
	)
}
