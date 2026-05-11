package api

import (
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/models"
	"github.com/shingo/server/repository"
	"github.com/shingo/server/utils/response"
	"go.mongodb.org/mongo-driver/bson"
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
		Email            string `json:"email"`
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		PropertyName     string `json:"property_name"`
		Phone            string `json:"phone"`
		QuestionnaireURL string `json:"questionnaire_url"`
		ThreadID         string `json:"thread_id"`
		CallID           string `json:"call_id"`
		Status           string `json:"status"`
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
		Email:            strings.TrimSpace(input.Email),
		FirstName:        strings.TrimSpace(input.FirstName),
		LastName:         strings.TrimSpace(input.LastName),
		PropertyName:     strings.TrimSpace(input.PropertyName),
		Phone:            strings.TrimSpace(input.Phone),
		QuestionnaireURL: strings.TrimSpace(input.QuestionnaireURL),
		ThreadID:         strings.TrimSpace(input.ThreadID),
		CallID:           strings.TrimSpace(input.CallID),
		Status:           status,
	})
	if err != nil {
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
	campaignID := c.Query("campaignId", c.Query("campaign_id", ""))

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

	contacts, total, err := h.repo.ListContacts(c.Context(), int64(page), int64(limit), search, status, campaignID)
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

func (h *ContactsHandler) UpdateContact(c fiber.Ctx) error {
	contactID := strings.TrimSpace(c.Params("id"))
	if contactID == "" {
		return response.BadRequest(c, "id is required")
	}

	var input map[string]any
	if err := c.Bind().Body(&input); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	updates := bson.M{}
	for key, value := range input {
		switch key {
		case "email":
			email, ok := value.(string)
			if !ok {
				return response.BadRequest(c, "email must be a string")
			}
			trimmed := strings.TrimSpace(email)
			if trimmed == "" {
				return response.BadRequest(c, "email must not be empty")
			}
			if _, err := mail.ParseAddress(trimmed); err != nil {
				return response.BadRequest(c, "email must be valid")
			}
			updates["email"] = trimmed
		case "first_name", "last_name", "property_name", "phone", "questionnaire_url", "thread_id", "call_id":
			strValue, ok := value.(string)
			if !ok {
				return response.BadRequest(c, key+" must be a string")
			}
			updates[key] = strings.TrimSpace(strValue)
		case "status":
			statusValue, ok := value.(string)
			if !ok {
				return response.BadRequest(c, "status must be a string")
			}
			normalizedStatus := models.ContactStatus(strings.ToLower(strings.TrimSpace(statusValue)))
			if normalizedStatus != models.ContactStatusActive &&
				normalizedStatus != models.ContactStatusUnsubscribed &&
				normalizedStatus != models.ContactStatusBounced {
				return response.BadRequest(c, "status must be one of active, unsubscribed or bounced")
			}
			updates["status"] = normalizedStatus
		}
	}

	if len(updates) == 0 {
		return response.BadRequest(c, "at least one updatable field is required")
	}

	contact, err := h.repo.UpdateContact(c.Context(), contactID, updates)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "contact not found")
		}

		return response.InternalError(c, "failed to update contact")
	}

	return response.Send(c, contact)
}

func (h *ContactsHandler) DeleteContact(c fiber.Ctx) error {
	contactID := strings.TrimSpace(c.Params("id"))
	if contactID == "" {
		return response.BadRequest(c, "id is required")
	}

	contact, err := h.repo.SoftDeleteContact(c.Context(), contactID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "contact not found")
		}

		return response.InternalError(c, "failed to delete contact")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code": http.StatusOK,
		"data": contact,
	})
}

func (h *ContactsHandler) ListDeletedContacts(c fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "25"))
	if err != nil {
		limit = 25
	}
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}

	contacts, total, err := h.repo.ListDeletedContacts(c.Context(), int64(page), int64(limit), search)
	if err != nil {
		return response.InternalError(c, "failed to fetch deleted contacts")
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

func (h *ContactsHandler) UpdateThreadID(c fiber.Ctx) error {
	contactID := strings.TrimSpace(c.Params("id"))
	threadID := strings.TrimSpace(c.Params("threadId"))

	if contactID == "" {
		return response.BadRequest(c, "id is required")
	}
	if threadID == "" {
		return response.BadRequest(c, "threadId is required")
	}

	contact, err := h.repo.UpdateThreadID(c.Context(), contactID, threadID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "contact not found")
		}

		return response.InternalError(c, "failed to update thread id")
	}

	return response.Send(c, contact)
}

func (h *ContactsHandler) UpdateCallID(c fiber.Ctx) error {
	contactID := strings.TrimSpace(c.Params("id"))
	callID := strings.TrimSpace(c.Params("callId"))

	if contactID == "" {
		return response.BadRequest(c, "id is required")
	}
	if callID == "" {
		return response.BadRequest(c, "callId is required")
	}

	contact, err := h.repo.UpdateCallID(c.Context(), contactID, callID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "contact not found")
		}

		return response.InternalError(c, "failed to update call id")
	}

	return response.Send(c, contact)
}

func (h *ContactsHandler) UpdateCampaignStatusByCallID(c fiber.Ctx) error {
	callID := strings.TrimSpace(c.Params("callId"))
	if callID == "" {
		return response.BadRequest(c, "callId is required")
	}

	var input struct {
		Status         string `json:"status"`
		NextCampaignID string `json:"next_campaign_id"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		return response.BadRequest(c, "status is required")
	}

	contact, err := h.repo.UpdateCampaignStatusByCallID(c.Context(), callID, status, strings.TrimSpace(input.NextCampaignID))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "contact not found")
		}

		return response.InternalError(c, "failed to update campaign status by call id")
	}

	return response.Send(c, contact)
}
