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
	apiGroup.Get("/deleted-contacts", contactsHandler.ListDeletedContacts)
	apiGroup.Get("/callback/:id", contactsHandler.MarkCallback)
	apiGroup.Get("/callback-handover", contactsHandler.CallbackHandover)
	apiGroup.Patch("/contact/:id", contactsHandler.UpdateContact)
	apiGroup.Patch("/contact/:id/thread-id/:threadId", contactsHandler.UpdateThreadID)
	apiGroup.Patch("/contact/:id/call-id/:callId", contactsHandler.UpdateCallID)
	apiGroup.Patch("/contact/call-id/:callId", contactsHandler.UpdateCampaignStatusByCallID)
	apiGroup.Delete("/contact/:id", contactsHandler.DeleteContact)
	apiGroup.Post("/disable", contactsHandler.DisableContact)
}

func RegisterCampaignRoutes(app fiber.Router, campaignsHandler *api.CampaignsHandler) {
	apiGroup := app.Group("/api/v1")

	apiGroup.Get("/campaigns", campaignsHandler.ListCampaigns)
	apiGroup.Get("/campaigns/:id", campaignsHandler.GetCampaign)
	apiGroup.Get("/campaigns/:id/contacts", campaignsHandler.ListCampaignContacts)
	apiGroup.Get("/handover/:id", campaignsHandler.Handover)
	apiGroup.Post("/campaigns/:id/contacts", campaignsHandler.AttachContacts)
	apiGroup.Post("/campaigns/:id/attach/:contactId", campaignsHandler.AttachContact)
	apiGroup.Patch("/campaigns/:id/contacts/:contactId/status", campaignsHandler.UpdateContactStatus)
	apiGroup.Delete("/campaigns/:id/contacts/:contactId", campaignsHandler.RemoveContact)
}
