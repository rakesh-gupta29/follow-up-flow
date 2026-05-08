package response

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/types"
)

// global 404 catchers

func APINotFound(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(types.Response{
		Code: fiber.StatusNotFound,
		Data: "API route not found",
	})
}

func WebNotFound(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(types.Response{
		Code: fiber.StatusNotFound,
		Data: "Web route not found",
	})
}

func Send(c fiber.Ctx, data any) error {
	return c.Status(http.StatusOK).JSON(types.Response{
		Code: http.StatusOK,
		Data: data,
	})
}

func SendError(c fiber.Ctx, code int, data any) error {
	return c.Status(code).JSON(types.Response{
		Code: code,
		Data: data,
	})
}

func BadRequest(c fiber.Ctx, msg string) error {
	return SendError(c, fiber.StatusBadRequest, msg)
}

func Unauthorized(c fiber.Ctx, msg string) error {
	return SendError(c, fiber.StatusUnauthorized, msg)
}

func Forbidden(c fiber.Ctx, msg string) error {
	return SendError(c, fiber.StatusForbidden, msg)
}

func NotFound(c fiber.Ctx, msg string) error {
	return SendError(c, fiber.StatusNotFound, msg)
}

func InternalError(c fiber.Ctx, msg string) error {
	return SendError(c, fiber.StatusInternalServerError, msg)
}
