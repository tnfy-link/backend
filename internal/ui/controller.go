package ui

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func Register(app *fiber.App) {
	app.Use(favicon.New(favicon.Config{
		File:       "/static/favicon.ico",
		FileSystem: http.FS(Static),
	}))

	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(Static),
		PathPrefix: "static",
		Browse:     false,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("templates/index", fiber.Map{})
	})
}
