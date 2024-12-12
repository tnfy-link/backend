package links

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/tnfy-link/server/internal/core/http/jsonify"
	validate "github.com/tnfy-link/server/internal/core/validator"
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
	if err != nil {
		c.l.Error("failed to get target", zap.Error(err), zap.String("id", id))
		return ctx.Redirect("/", fiber.StatusTemporaryRedirect)
	}

	return ctx.Redirect(target, fiber.StatusTemporaryRedirect)
}

func (c *controller) post(ctx *fiber.Ctx) error {
	req := PostLinksRequest{}
	if err := c.bodyParserValidator(ctx, &req); err != nil {
		return err
	}

	link, err := c.s.Create(ctx.Context(), req.Link)
	if err != nil {
		c.l.Error("failed to create link", zap.Error(err), zap.String("link", req.Link.TargetURL))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create link: %s", err.Error()))
	}

	link.URL = fmt.Sprintf("%s/%s", c.hostname, link.ID)

	return ctx.JSON(
		PostLinksResponse{
			Link: link,
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
	app.Get("/:id", c.redirect)

	api := app.Group("/api/v1", jsonify.New())
	api.Post(
		"/links",
		NewLimiter(),
		c.post,
	)

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
