package main

import (
	"context"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tnfy-link/server/internal/core/api"
	"github.com/tnfy-link/server/internal/core/http"
	"github.com/tnfy-link/server/internal/core/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	fx.New(
		logger.Module,
		http.Module,
		api.Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			logOption := fxevent.ZapLogger{Logger: logger}
			logOption.UseLogLevel(zapcore.DebugLevel)
			return &logOption
		}),
		fx.Invoke(func(api hureg.APIGen, logger *zap.Logger) {
			hureg.Register(
				api,
				huma.Operation{
					Method: "GET",
					Path:   "/",
				},
				func(ctx context.Context, i *struct{}) (*struct{}, error) {
					return nil, nil
				},
			)
		}),
	).
		Run()
}
