package links

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/tnfy-link/server/internal/core/http/jsonify"
	"go.uber.org/zap"
)

type controller struct {
	r *repository
	l *zap.Logger

	hostname string
}

func (c *controller) redirect(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	target, err := c.r.GetTarget(ctx.Context(), id)
	if err != nil {
		c.l.Error("failed to get target", zap.Error(err), zap.String("id", id))
		return ctx.Redirect("/", fiber.StatusTemporaryRedirect)
	}

	return ctx.Redirect(target, fiber.StatusTemporaryRedirect)
}

func (c *controller) post(ctx *fiber.Ctx) error {
	req := PostLinksRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to parse request: %s", err.Error()))
	}

	link, err := c.r.Create(ctx.Context(), req.Link)
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

func newController(r *repository, l *zap.Logger, c Config) *controller {
	return &controller{
		r: r,
		l: l,

		hostname: c.Hostname,
	}
}

func Register(app *fiber.App, redis *redis.Client, config Config, log *zap.Logger) error {
	repository := newRepository(redis, config)
	controller := newController(repository, log, config)

	app.Get("/:id", controller.redirect)

	api := app.Group("/api/v1", jsonify.New())
	api.Post(
		"/links",
		NewLimiter(),
		controller.post,
	)

	return nil
}
