package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/utils/response"
)

type AppHandler struct{}

func NewAppHandler() *AppHandler {
	return &AppHandler{}
}

func (h *AppHandler) HealthCheck(c fiber.Ctx) error {
	return response.Send(c, "API: health check")
}
