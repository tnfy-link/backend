package http

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func New(config Config, views fiber.Views, logger *zap.Logger) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          config.Proxies,
		ProxyHeader:             config.ProxyHeader,
		EnableIPValidation:      true,
		Views:                   views,
	})
	app.Use(fiberzap.New(fiberzap.Config{
		SkipBody: func(c *fiber.Ctx) bool {
			return c.Response().StatusCode() < 400
		},
		Logger: logger,
		Fields: []string{"latency", "status", "method", "url", "ip", "ua", "body", "error"},
	}))
	app.Use(recover.New())

	return app, nil
}
