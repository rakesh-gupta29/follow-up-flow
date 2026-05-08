package apiRoutes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/handlers/api"
	"github.com/shingo/server/utils/response"
)

func RegisterAdminRoutes(app fiber.Router, authHandler *api.AuthHandler) {
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/login", authHandler.Login)
}

func RegisterAppRoutes(app fiber.Router, appHandler *api.AppHandler) {
	apiGroup := app.Group("/api/v1")
	apiGroup.Get("/health", appHandler.HealthCheck)

	apiGroup.Use(response.APINotFound)
}
