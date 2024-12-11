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
		EnableIPValidation:      true,
		EnableTrustedProxyCheck: true,
		ErrorHandler:            ErrorHandler,
		ProxyHeader:             config.ProxyHeader,
		TrustedProxies:          config.Proxies,
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

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Send json error
	return c.Status(code).JSON(ErrorResponse{Error: Error{Message: err.Error()}})
}
