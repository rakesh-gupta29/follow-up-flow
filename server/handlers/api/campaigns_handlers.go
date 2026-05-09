package api

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/shingo/server/models"
	"github.com/shingo/server/repository"
	"github.com/shingo/server/utils/response"
	"go.mongodb.org/mongo-driver/mongo"
)

type CampaignsHandler struct {
	repo *repository.CampaignsRepository
}

func NewCampaignsHandler(repo *repository.CampaignsRepository) *CampaignsHandler {
	return &CampaignsHandler{repo: repo}
}

func (h *CampaignsHandler) ListCampaigns(c fiber.Ctx) error {
	campaigns, err := h.repo.ListCampaigns(c.Context())
	if err != nil {
		return response.InternalError(c, "failed to fetch campaigns")
	}

	return response.Send(c, campaigns)
}

func (h *CampaignsHandler) GetCampaign(c fiber.Ctx) error {
	campaign, err := h.repo.GetCampaign(c.Context(), strings.TrimSpace(c.Params("id")))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "campaign not found")
		}

		return response.InternalError(c, "failed to fetch campaign")
	}

	return response.Send(c, campaign)
}

func (h *CampaignsHandler) ListCampaignContacts(c fiber.Ctx) error {
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

	campaignID := strings.TrimSpace(c.Params("id"))
	contacts, total, err := h.repo.ListCampaignContacts(c.Context(), campaignID, int64(page), int64(limit), search, status)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "campaign not found")
		}

		return response.InternalError(c, "failed to fetch campaign contacts")
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

func (h *CampaignsHandler) AttachContact(c fiber.Ctx) error {
	campaignID := strings.TrimSpace(c.Params("id"))
	contactID := strings.TrimSpace(c.Params("contactId"))

	if campaignID == "" || contactID == "" {
		return response.BadRequest(c, "campaign id and contact id are required")
	}

	contact, err := h.repo.AttachContact(c.Context(), campaignID, contactID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "campaign or contact not found")
		}

		return response.InternalError(c, "failed to attach contact to campaign")
	}

	return response.Send(c, contact)
}

func (h *CampaignsHandler) AttachContacts(c fiber.Ctx) error {
	campaignID := strings.TrimSpace(c.Params("id"))
	if campaignID == "" {
		return response.BadRequest(c, "campaign id is required")
	}

	var input struct {
		ContactIDs []string `json:"contact_ids"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	if len(input.ContactIDs) == 0 {
		return response.BadRequest(c, "contact_ids is required")
	}

	modifiedCount, err := h.repo.AttachContacts(c.Context(), campaignID, input.ContactIDs)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "campaign not found")
		}

		return response.InternalError(c, "failed to attach contacts to campaign")
	}

	return response.Send(c, fiber.Map{
		"campaign_id":    campaignID,
		"attached_count": modifiedCount,
		"contact_ids":    input.ContactIDs,
	})
}

func (h *CampaignsHandler) RemoveContact(c fiber.Ctx) error {
	campaignID := strings.TrimSpace(c.Params("id"))
	contactID := strings.TrimSpace(c.Params("contactId"))

	if campaignID == "" || contactID == "" {
		return response.BadRequest(c, "campaign id and contact id are required")
	}

	contact, err := h.repo.RemoveContact(c.Context(), campaignID, contactID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "campaign or contact mapping not found")
		}

		return response.InternalError(c, "failed to remove contact from campaign")
	}

	return response.Send(c, contact)
}

func (h *CampaignsHandler) UpdateContactStatus(c fiber.Ctx) error {
	campaignID := strings.TrimSpace(c.Params("id"))
	contactID := strings.TrimSpace(c.Params("contactId"))

	if campaignID == "" || contactID == "" {
		return response.BadRequest(c, "campaign id and contact id are required")
	}

	var input struct {
		Status         string `json:"status"`
		NextCampaignID string `json:"next_campaign_id"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.BadRequest(c, "Invalid input")
	}

	status := models.CampaignContactStatus(strings.TrimSpace(input.Status))
	if status != models.CampaignContactStatusQueued &&
		status != models.CampaignContactStatusInProgress &&
		status != models.CampaignContactStatusSuccess &&
		status != models.CampaignContactStatusFailed {
		return response.BadRequest(c, "status must be one of queued, in_progress, success or failed")
	}

	contact, err := h.repo.UpdateCampaignContactStatus(c.Context(), campaignID, contactID, status, strings.TrimSpace(input.NextCampaignID))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "campaign contact not found")
		}

		return response.InternalError(c, "failed to update campaign contact status")
	}

	return response.Send(c, contact)
}

func (h *CampaignsHandler) Handover(c fiber.Ctx) error {
	campaignID := strings.TrimSpace(c.Params("id"))
	if campaignID == "" {
		return response.BadRequest(c, "campaign id is required")
	}

	contacts, err := h.repo.GetQueuedCampaignContacts(c.Context(), campaignID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFound(c, "campaign not found")
		}

		return response.InternalError(c, "failed to fetch queued contacts")
	}

	return response.Send(c, contacts)
}
