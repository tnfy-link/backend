package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/tnfy-link/backend/internal/api/limiter"
	"github.com/tnfy-link/backend/internal/api/param"
	"github.com/tnfy-link/backend/internal/links"
	"github.com/tnfy-link/client-go/api"
	"github.com/tnfy-link/core/handler"
	"go.uber.org/zap"
)

type Links struct {
	handler.Base

	s *links.Service
}

func (c *Links) get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	link, err := c.s.Get(ctx.Context(), id)
	if err == links.ErrLinkNotFound {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		c.Logger.Error("failed to get link", zap.Error(err), zap.String("id", id))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(
		api.GetLinkResponse{
			Link: link,
		},
	)
}

func (c *Links) post(ctx *fiber.Ctx) error {
	req := api.PostLinksRequest{}
	if err := c.BodyParserValidator(ctx, &req); err != nil {
		return err
	}

	link, err := c.s.Create(ctx.Context(), req.Link)
	if verr := links.AsValidationError(err); verr != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to create link: %s", verr.Error()))
	}
	if err != nil {
		c.Logger.Error("failed to create link", zap.Error(err), zap.String("link", req.Link.TargetURL))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create link: %s", err.Error()))
	}

	return ctx.JSON(
		api.PostLinksResponse{
			Link: link,
		},
	)
}

func (c *Links) stats(ctx *fiber.Ctx) error {
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

func (c *Links) Register(router fiber.Router) {
	idValidator := param.NewValidator("id", c.s.ValidateID)

	router.Get("/:id", idValidator, c.get)
	router.Post(
		"/",
		limiter.New(1),
		c.post,
	)
	router.Get("/:id/stats", idValidator, c.stats)
}

func NewLinks(s *links.Service, v *validator.Validate, l *zap.Logger) *Links {
	return &Links{
		Base: handler.Base{
			Validator: v,
			Logger:    l,
		},
		s: s,
	}
}
