package api

import (
	"net/mail"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/models"
	"github.com/shingo/server/repository"
	"github.com/shingo/server/utils/response"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactsHandler struct {
	repo *repository.ContactsRepository
}

func NewContactsHandler(repo *repository.ContactsRepository) *ContactsHandler {
	return &ContactsHandler{repo: repo}
}

func (h *ContactsHandler) AddContact(c fiber.Ctx) error {
	var input struct {
		Email     string            `json:"email"`
		FirstName string            `json:"first_name"`
		LastName  string            `json:"last_name"`
		Phone     string            `json:"phone"`
		Company   string            `json:"company"`
		Tags      []string          `json:"tags"`
		Status    string            `json:"status"`
		Meta      map[string]string `json:"meta"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.LastName) == "" {
		return response.BadRequest(c, "email, first_name and last_name are required")
	}
	if _, err := mail.ParseAddress(strings.TrimSpace(input.Email)); err != nil {
		return response.BadRequest(c, "email must be valid")
	}

	status := models.ContactStatus(strings.TrimSpace(input.Status))
	if status == "" {
		status = models.ContactStatusActive
	}

	if status != models.ContactStatusActive && status != models.ContactStatusUnsubscribed && status != models.ContactStatusBounced {
		return response.BadRequest(c, "status must be one of active, unsubscribed or bounced")
	}

	contact, err := h.repo.CreateContact(c.Context(), models.Contact{
		Email:     strings.TrimSpace(input.Email),
		FirstName: strings.TrimSpace(input.FirstName),
		LastName:  strings.TrimSpace(input.LastName),
		Phone:     strings.TrimSpace(input.Phone),
		Company:   strings.TrimSpace(input.Company),
		Tags:      input.Tags,
		Status:    status,
		Meta:      input.Meta,
	})
	if err != nil {
		if repository.IsDuplicateKeyError(err) {
			return response.BadRequest(c, "contact with this email already exists")
		}

		return response.InternalError(c, "failed to create contact")
	}

	return response.Send(c, contact)
}

func (h *ContactsHandler) ListContacts(c fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "25"))
	if err != nil {
		limit = 25
	}
	search := c.Query("search", "")
	status := c.Query("status", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}

	if status != "" {
		normalizedStatus := strings.ToLower(strings.TrimSpace(status))
		if normalizedStatus != string(models.ContactStatusActive) &&
			normalizedStatus != string(models.ContactStatusUnsubscribed) &&
			normalizedStatus != string(models.ContactStatusBounced) {
			return response.BadRequest(c, "status must be one of active, unsubscribed or bounced")
		}
		status = normalizedStatus
	}

	contacts, total, err := h.repo.ListContacts(c.Context(), int64(page), int64(limit), search, status)
	if err != nil {
		return response.InternalError(c, "failed to fetch contacts")
	}

	return response.Send(c, fiber.Map{
		"items": contacts,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *ContactsHandler) DisableContact(c fiber.Ctx) error {
	var input struct {
		ID string `json:"id"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	if strings.TrimSpace(input.ID) == "" {
		return response.BadRequest(c, "id is required")
	}

	contact, err := h.repo.DisableContact(c.Context(), strings.TrimSpace(input.ID))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "contact not found")
		}

		return response.InternalError(c, "failed to disable contact")
	}

	return response.Send(c, contact)
}
