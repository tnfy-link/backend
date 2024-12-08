package main

import (
	"github.com/tnfy-link/server/internal/config"
	"github.com/tnfy-link/server/internal/core/http"
	"github.com/tnfy-link/server/internal/core/logger"
	"github.com/tnfy-link/server/internal/core/redis"
	"github.com/tnfy-link/server/internal/links"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	fx.New(
		logger.Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			logOption := fxevent.ZapLogger{Logger: logger}
			logOption.UseLogLevel(zapcore.DebugLevel)
			return &logOption
		}),
		http.Module,
		redis.Module,

		config.Module,
		links.Module,
	).
		Run()
}
