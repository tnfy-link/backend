package links

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tnfy-link/backend/internal/core/http/jsonify"
	validate "github.com/tnfy-link/backend/internal/core/validator"
	"go.uber.org/zap"
)

type controller struct {
	s *service

	v *validator.Validate
	l *zap.Logger

	hostname string
}

func (c *controller) redirect(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	target, err := c.s.GetTarget(ctx.Context(), id)
	if err == ErrInvalidID {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if err == ErrLinkNotFound {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		c.l.Error("failed to get target", zap.Error(err), zap.String("id", id))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	go func(id, query string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		if err := c.s.RegisterStats(ctx, id, query); err != nil {
			c.l.Error("failed to register stats", zap.Error(err), zap.String("id", id), zap.String("query", query))
		}
	}(strings.Clone(id), strings.Clone(ctx.Context().QueryArgs().String()))

	return ctx.Redirect(target, fiber.StatusTemporaryRedirect)
}

func (c *controller) get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	link, err := c.s.Get(ctx.Context(), id)
	if err == ErrInvalidID {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if err == ErrLinkNotFound {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		c.l.Error("failed to get link", zap.Error(err), zap.String("id", id))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	link.URL = c.hostname + "/" + link.ID

	return ctx.JSON(
		GetLinksResponse{
			Link: link,
		},
	)
}

func (c *controller) post(ctx *fiber.Ctx) error {
	req := PostLinksRequest{}
	if err := c.bodyParserValidator(ctx, &req); err != nil {
		return err
	}

	link, err := c.s.Create(ctx.Context(), req.Link)
	if err := AsValidationError(err); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to create link: %s", err.Error()))
	}
	if err != nil {
		c.l.Error("failed to create link", zap.Error(err), zap.String("link", req.Link.TargetURL))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create link: %s", err.Error()))
	}

	link.URL = c.hostname + "/" + link.ID

	return ctx.JSON(
		PostLinksResponse{
			Link: link,
		},
	)
}

func (c *controller) stats(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	stats, err := c.s.GetStats(ctx.Context(), id)
	if err != nil {
		c.l.Error("failed to get stats", zap.Error(err), zap.String("id", id))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get stats: %s", err.Error()))
	}

	return ctx.JSON(
		GetStatsResponse{
			Stats: stats,
		},
	)
}

func (c *controller) bodyParserValidator(ctx *fiber.Ctx, out interface{}) error {
	if err := ctx.BodyParser(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to parse request: %s", err.Error()))
	}

	if err := c.v.Struct(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to validate request: %s", err.Error()))
	}

	if v, ok := out.(validate.Validatable); ok {
		if err := v.Validate(); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to validate request: %s", err.Error()))
		}
	}

	return nil
}

func (c *controller) Register(app *fiber.App) {
	api := app.Group(
		"/api/v1",
		cors.New(cors.Config{
			AllowOrigins: c.hostname,
		}),
		jsonify.New(),
	)
	api.Get("/links/:id", c.get)
	api.Post(
		"/links",
		NewLimiter(),
		c.post,
	)
	api.Get("/links/:id/stats", c.stats)
	api.Use(func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusNotFound)
	})

	app.Get("/:id", c.redirect)
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/", fiber.StatusTemporaryRedirect)
	})
}

func newController(s *service, v *validator.Validate, l *zap.Logger, c Config) *controller {
	return &controller{
		s: s,

		v: v,
		l: l,

		hostname: c.Hostname,
	}
}
