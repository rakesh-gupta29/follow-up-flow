package webRoutes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/utils/response"
)

func RegisterRoutes(app fiber.Router) {
	webGroup := app.Group("/")

	webGroup.Use(response.WebNotFound)
}
