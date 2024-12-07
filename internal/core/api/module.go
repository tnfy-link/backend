package api

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"api",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("api")
	}),
	fx.Provide(New),
)
