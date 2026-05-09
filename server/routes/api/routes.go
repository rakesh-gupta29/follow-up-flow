package apiRoutes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/handlers/api"
	"github.com/shingo/server/utils/response"
)

func RegisterAdminRoutes(app fiber.Router, authHandler *api.AuthHandler) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", authHandler.Login)
	auth.Get("/me", authHandler.Me)
}

func RegisterAppRoutes(app fiber.Router, appHandler *api.AppHandler) {
	apiGroup := app.Group("/api/v1")
	apiGroup.Get("/health", appHandler.HealthCheck)

	apiGroup.Use(response.APINotFound)
}

func RegisterContactsRoutes(app fiber.Router, contactsHandler *api.ContactsHandler) {
	apiGroup := app.Group("/api/v1")
	apiGroup.Post("/add-contact", contactsHandler.AddContact)
	apiGroup.Get("/contacts", contactsHandler.ListContacts)
	apiGroup.Post("/disable", contactsHandler.DisableContact)
}
