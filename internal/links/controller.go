package links

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type controller struct {
	r *repository
	l *zap.Logger
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

	return ctx.JSON(
		PostLinksResponse{
			Link: link,
		},
	)
}

func newController(r *repository, l *zap.Logger) *controller {
	return &controller{
		r: r,
		l: l,
	}
}

func Register(app *fiber.App, redis *redis.Client, log *zap.Logger) error {
	repository := newRepository(redis)
	controller := newController(repository, log)

	api := app.Group("/api")

	api.Post("/links", controller.post)

	return nil
}
