package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnfy-link/core/http"
	"github.com/tnfy-link/core/http/jsonify"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"api",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("api")
	}),
	fx.Provide(http.NewJSONErrorHandler),
	fx.Provide(func(log *zap.Logger) http.Options {
		return *(&http.Options{}).WithErrorHandler(http.NewJSONErrorHandler(log))
	}),
	fx.Provide(NewLinks),
	fx.Invoke(func(app *fiber.App, l *Links) {
		api := app.Group(
			"/api/v1",
			jsonify.New(),
		)

		l.Register(api.Group("/links"))

		api.Use(func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusNotFound)
		})
	}),
)
