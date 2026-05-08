// internal/handlers/api/auth_handler.go
package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/repository"
	"github.com/shingo/server/utils/response"
	"go.mongodb.org/mongo-driver/mongo"
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

	user, err := h.repo.FindByEmail(c.Context(), input.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "User not found")
		}
		return response.InternalError(c, "Database error")
	}

	// Simple password check (In production, use bcrypt.CompareHashAndPassword)
	if user.Password != input.Password {
		return response.BadRequest(c, "Incorrect password")
	}

	return response.Send(c, fiber.Map{"message": "Login successful", "email": user.Email})
}
