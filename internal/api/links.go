package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/tnfy-link/backend/internal/api/limiter"
	"github.com/tnfy-link/backend/internal/api/param"
	"github.com/tnfy-link/backend/internal/links"
	"github.com/tnfy-link/backend/internal/stats"
	"github.com/tnfy-link/client-go/api"
	"github.com/tnfy-link/core/handler"
	"go.uber.org/zap"
)

type Links struct {
	handler.Base

	links *links.Service
	stats *stats.Service
}

//	@Summary		Get link metadata
//	@Description	Get link metadata by ID
//	@Tags			Links
//	@Produce		json
//	@Param			id	path		string	true	"Link ID"
//	@Success		200	{object}	api.GetLinkResponse
//	@Failure		400	{object}	http.JSONErrorResponse
//	@Failure		404	{object}	http.JSONErrorResponse
//	@Failure		500	{object}	http.JSONErrorResponse
//	@Router			/links/{id} [get]
//
// Get Link
func (c *Links) get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	link, err := c.links.Get(ctx.Context(), id)
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

//	@Summary		Shorten link
//	@Description	Shorten link
//	@Tags			Links
//	@Accept			json
//	@Produce		json
//	@Param			link	body		api.PostLinksRequest	true	"Link"
//	@Success		201		{object}	api.PostLinksResponse
//	@Failure		400		{object}	http.JSONErrorResponse "Bad Request"
//	@Failure		429		{object}	http.JSONErrorResponse "Too many requests"
//	@Failure		500		{object}	http.JSONErrorResponse
//	@Router			/links [post]
//
// Shorten link
func (c *Links) post(ctx *fiber.Ctx) error {
	req := api.PostLinksRequest{}
	if err := c.BodyParserValidator(ctx, &req); err != nil {
		return err
	}

	link, err := c.links.Create(ctx.Context(), links.NewNewLink(req.Link.TargetURL))
	if verr := links.AsValidationError(err); verr != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to create link: %s", verr.Error()))
	}
	if err != nil {
		c.Logger.Error("failed to create link", zap.Error(err), zap.String("link", req.Link.TargetURL))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create link: %s", err.Error()))
	}

	return ctx.
		Status(fiber.StatusCreated).
		JSON(
			api.PostLinksResponse{
				Link: link,
			},
		)
}

//	@Summary		Get link stats
//	@Description	Get link stats by ID
//	@Tags			Links
//	@Produce		json
//	@Param			id	path		string	true	"Link ID"
//	@Success		200	{object}	api.GetStatsResponse
//	@Failure		400	{object}	http.JSONErrorResponse
//	@Failure		404	{object}	http.JSONErrorResponse
//	@Failure		500	{object}	http.JSONErrorResponse
//	@Router			/links/{id}/stats [get]
//
// Get Statistics
func (c *Links) getStats(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	s, err := c.stats.Get(ctx.Context(), id)
	if err == stats.ErrNotFound {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		c.Logger.Error("failed to get stats", zap.Error(err), zap.String("id", id))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get link stats: %s", err.Error()))
	}

	return ctx.JSON(
		api.GetStatsResponse{
			Stats: s,
		},
	)
}

func (c *Links) Register(router fiber.Router) {
	idValidator := param.NewValidator("id", c.links.ValidateID)

	router.Get("/:id", idValidator, c.get)
	router.Post(
		"/",
		limiter.New(1),
		c.post,
	)
	router.Get("/:id/stats", idValidator, c.getStats)
}

func NewLinks(links *links.Service, stats *stats.Service, v *validator.Validate, l *zap.Logger) *Links {
	switch {
	case links == nil:
		panic("links service is required")
	case stats == nil:
		panic("stats service is required")
	case v == nil:
		panic("validator is required")
	case l == nil:
		panic("logger is required")

	}

	return &Links{
		Base: handler.Base{
			Validator: v,
			Logger:    l,
		},
		links: links,
		stats: stats,
	}
}
