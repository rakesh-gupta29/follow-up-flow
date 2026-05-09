// internal/handlers/api/auth_handler.go
package api

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/repository"
	"github.com/shingo/server/utils/response"
)

type AuthHandler struct {
	repo *repository.AdminRepository
}

func NewAuthHandler(repo *repository.AdminRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	user, err := h.repo.GetAdmin(c.Context(), input.Email)
	fmt.Println("admin foiudn", user)
	if err != nil {
		return response.BadRequest(c, "Unauthorized: Invalid credentials")
	}
	if user.Password != input.Password {
		fmt.Println("invalid password", user.Password, input.Password)
		return response.BadRequest(c, "Unauthorized: Invalid credentials")
	}

	// Simple Token for NudgeBuddy
	return response.Send(c, fiber.Map{
		"access_token": "nudgebuddy_secret_token_123",
		"user":         user.Email,
	})
}

func (h *AuthHandler) Me(c fiber.Ctx) error {
	token := c.Get("Authorization")

	if token != "Bearer nudgebuddy_secret_token_123" {
		return response.BadRequest(c, "Invalid or expired token")
	}

	return response.Send(c, fiber.Map{
		"email": "admin",
		"role":  "owner",
	})
}
