package ui

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func NewViews() fiber.Views {
	return html.NewFileSystem(http.FS(Templates), ".html")
}
