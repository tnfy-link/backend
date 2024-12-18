package links

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/tnfy-link/backend/internal/api/param"
	"github.com/tnfy-link/backend/internal/core/handler"
	"github.com/tnfy-link/client-go/api"
	"go.uber.org/zap"
)

type controller struct {
	handler.Base

	s *Service

	hostname string
}

func (c *controller) redirect(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	target, err := c.s.GetTarget(ctx.Context(), id)
	if err == ErrLinkNotFound {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		c.Logger.Error("failed to get target", zap.Error(err), zap.String("id", id))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	go func(id, query string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		if err := c.s.RegisterStats(ctx, id, query); err != nil {
			c.Logger.Error("failed to register stats", zap.Error(err), zap.String("id", id), zap.String("query", query))
		}
	}(strings.Clone(id), strings.Clone(ctx.Context().QueryArgs().String()))

	return ctx.Redirect(target, fiber.StatusTemporaryRedirect)
}

func (c *controller) stats(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	stats, err := c.s.GetStats(ctx.Context(), id)
	if err != nil {
		c.Logger.Error("failed to get stats", zap.Error(err), zap.String("id", id))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get stats: %s", err.Error()))
	}

	return ctx.JSON(
		api.GetStatsResponse{
			Stats: stats,
		},
	)
}

func (c *controller) Register(app *fiber.App) {
	idValidator := param.NewValidator("id", c.s.ValidateID)

	// api := app.Group(
	// 	"/api/v1",
	// 	cors.New(cors.Config{
	// 		AllowOrigins: c.hostname,
	// 	}),
	// 	jsonify.New(),
	// )
	// api.Get("/links/:id/stats", idValidator, c.stats)
	// api.Use(func(ctx *fiber.Ctx) error {
	// 	return ctx.SendStatus(fiber.StatusNotFound)
	// })

	app.Get("/:id", idValidator, c.redirect)
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/", fiber.StatusTemporaryRedirect)
	})
}

func newController(s *Service, v *validator.Validate, l *zap.Logger, c Config) *controller {
	return &controller{
		Base: handler.Base{
			Validator: v,
			Logger:    l,
		},
		s:        s,
		hostname: c.Hostname,
	}
}
