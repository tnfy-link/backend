package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/tnfy-link/core/http"
	"github.com/tnfy-link/core/http/jsonify"
	"go.uber.org/fx"
	"go.uber.org/zap"

	apidoc "github.com/tnfy-link/backend/api"
	"github.com/tnfy-link/backend/internal/version"
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
	fx.Invoke(func(app *fiber.App, l *Links, config Config) {
		api := app.Group("/api/v1")

		apidoc.SwaggerInfo.Version = version.AppVersion
		api.Get("/docs/*", swagger.New(swagger.Config{
			Layout: "BaseLayout",
		}))

		if config.CORSAllowOrigins != "" {
			api.Use(cors.New(cors.Config{
				AllowOrigins:     config.CORSAllowOrigins,
				AllowCredentials: true,
				MaxAge:           86400,
			}))
		}

		api.Use(jsonify.New())

		l.Register(api.Group("/links"))

		api.Use(func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusNotFound)
		})
	}),
)
