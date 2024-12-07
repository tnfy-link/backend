package api

import (
	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App) hureg.APIGen {
	cfg := huma.DefaultConfig("Tnfy API", "1.0.0")
	// cfg.DocsPath = ""
	// cfg.OpenAPIPath = "api"
	// cfg.SchemasPath = "/schemas"

	// api := humafiber.NewWithGroup(app, app.Group("/api"), cfg)
	api := humafiber.New(app, cfg)
	return hureg.NewAPIGen(api).AddBasePath("/api")
	// return humafiber.New(app, huma.DefaultConfig("Tnfy API", "1.0.0"))
}
