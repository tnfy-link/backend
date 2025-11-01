package internal

import (
	"github.com/tnfy-link/backend/internal/api"
	"github.com/tnfy-link/backend/internal/config"
	"github.com/tnfy-link/backend/internal/id"
	"github.com/tnfy-link/backend/internal/links"
	"github.com/tnfy-link/backend/internal/stats"
	"github.com/tnfy-link/core/http"
	"github.com/tnfy-link/core/logger"
	"github.com/tnfy-link/core/redis"
	"github.com/tnfy-link/core/validator"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run() {
	fx.New(
		logger.Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			logOption := fxevent.ZapLogger{Logger: logger}
			logOption.UseLogLevel(zapcore.DebugLevel)
			return &logOption
		}),
		http.Module,
		redis.Module,
		validator.Module,

		config.Module(),
		api.Module(),
		id.Module(),
		links.Module(),
		stats.Module(),
	).
		Run()
}
