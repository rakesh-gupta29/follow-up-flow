package web

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/utils/response"
)

func HealthCheck(c fiber.Ctx) error {
	return response.Send(c, "Web: health check")
}
