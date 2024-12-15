package internal

import (
	"github.com/tnfy-link/backend/internal/config"
	"github.com/tnfy-link/backend/internal/core/http"
	"github.com/tnfy-link/backend/internal/core/logger"
	"github.com/tnfy-link/backend/internal/core/redis"
	"github.com/tnfy-link/backend/internal/core/validator"
	"github.com/tnfy-link/backend/internal/id"
	"github.com/tnfy-link/backend/internal/links"
	"github.com/tnfy-link/backend/internal/ui"
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

		config.Module,
		ui.Module,
		id.Module,
		links.Module,
	).
		Run()
}
