package links

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"links",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("api")
	}),
	fx.Provide(newRepository, fx.Private),
	fx.Provide(newService, fx.Private),
	fx.Provide(newController, fx.Private),
	fx.Invoke(func(c *controller, app *fiber.App) {
		c.Register(app)
	}),
)
