package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"http",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("http")
	}),
	fx.Provide(New),
	fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, logger *zap.Logger) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					_ = app.Listen(":3000")
				}()
				logger.Info("server started")

				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Warn("shutting down server")
				_ = app.ShutdownWithContext(ctx)
				return nil
			},
		})
	}),
)
