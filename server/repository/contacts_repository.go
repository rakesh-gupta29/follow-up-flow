package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shingo/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const appDatabaseName = "nudgebuddy_db"

type ContactsRepository struct {
	client              *mongo.Client
	collection          *mongo.Collection
	campaignsCollection *mongo.Collection
}

func NewContactsRepository(db *mongo.Client) *ContactsRepository {
	return &ContactsRepository{
		client:              db,
		collection:          db.Database(appDatabaseName).Collection("contacts"),
		campaignsCollection: db.Database(appDatabaseName).Collection("campaigns"),
	}
}

func (r *ContactsRepository) EnsureCollection(ctx context.Context) error {
	db := r.client.Database(appDatabaseName)

	collections, err := db.ListCollectionNames(ctx, bson.M{"name": "contacts"})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		if err := db.CreateCollection(ctx, "contacts"); err != nil {
			// If another process created it first, Mongo will return a namespace exists error.
			if !strings.Contains(strings.ToLower(err.Error()), "exists") {
				return err
			}
		}
	}

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "campaign_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "campaign_ids", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "campaign_memberships.campaign_id", Value: 1}},
		},
	}

	_, err = r.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

func (r *ContactsRepository) CreateContact(ctx context.Context, contact models.Contact) (*models.Contact, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	contact.ID = uuid.NewString()
	contact.Status = normalizeStatus(contact.Status)
	contact.CreatedAt = now
	contact.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *ContactsRepository) ListContacts(ctx context.Context, page, limit int64, search string, status string, campaignID string) ([]models.ContactListItem, int64, error) {
	filter := bson.M{}

	if trimmedSearch := strings.TrimSpace(search); trimmedSearch != "" {
		pattern := regexp.QuoteMeta(trimmedSearch)
		filter["$or"] = []bson.M{
			{"email": bson.M{"$regex": pattern, "$options": "i"}},
			{"first_name": bson.M{"$regex": pattern, "$options": "i"}},
			{"last_name": bson.M{"$regex": pattern, "$options": "i"}},
			{"company": bson.M{"$regex": pattern, "$options": "i"}},
		}
	}

	if normalizedStatus := normalizeStatus(models.ContactStatus(status)); normalizedStatus != "" {
		filter["status"] = normalizedStatus
	}
	if strings.TrimSpace(campaignID) != "" {
		trimmedCampaignID := strings.TrimSpace(campaignID)
		filter["$and"] = []bson.M{
			{
				"$or": []bson.M{
					{"campaign_id": trimmedCampaignID},
					{"campaign_ids": trimmedCampaignID},
					{"campaign_memberships.campaign_id": trimmedCampaignID},
				},
			},
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	contacts := make([]models.ContactListItem, 0)
	for cursor.Next(ctx) {
		var contact models.Contact
		if err := cursor.Decode(&contact); err != nil {
			return nil, 0, err
		}

		memberships := normalizeCampaignMemberships(contact)
		for index := range memberships {
			var campaignDoc models.Campaign
			err := r.campaignsCollection.FindOne(ctx, bson.M{"id": memberships[index].CampaignID}).Decode(&campaignDoc)
			if err != nil && err != mongo.ErrNoDocuments {
				return nil, 0, err
			}
			if err == nil {
				memberships[index].Campaign = &campaignDoc
			}
		}
		contact.CampaignMemberships = memberships

		campaigns := make([]models.ContactCampaign, 0, len(memberships))
		for _, membership := range memberships {
			var campaignDoc models.Campaign
			err := r.campaignsCollection.FindOne(ctx, bson.M{"id": membership.CampaignID}).Decode(&campaignDoc)
			if err != nil && err != mongo.ErrNoDocuments {
				return nil, 0, err
			}
			if err == nil {
				campaigns = append(campaigns, models.ContactCampaign{
					ID:     campaignDoc.ID,
					Name:   campaignDoc.Name,
					Status: membership.Status,
				})
			}
		}

		var campaign *models.ContactCampaign
		if len(campaigns) > 0 {
			campaign = &campaigns[0]
		}

		contacts = append(contacts, models.ContactListItem{
			Contact:   contact,
			Campaign:  campaign,
			Campaigns: campaigns,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (r *ContactsRepository) DisableContact(ctx context.Context, id string) (*models.Contact, error) {
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     models.ContactStatusUnsubscribed,
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		},
	}

	result := r.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		return nil, result.Err()
	}

	var contact models.Contact
	if err := result.Decode(&contact); err != nil {
		return nil, err
	}

	return &contact, nil
}

func normalizeStatus(status models.ContactStatus) models.ContactStatus {
	switch strings.ToLower(string(status)) {
	case string(models.ContactStatusActive):
		return models.ContactStatusActive
	case string(models.ContactStatusUnsubscribed):
		return models.ContactStatusUnsubscribed
	case string(models.ContactStatusBounced):
		return models.ContactStatusBounced
	default:
		return ""
	}
}

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	return mongo.IsDuplicateKeyError(err) || strings.Contains(strings.ToLower(fmt.Sprint(err)), "e11000")
}

func normalizeCampaignIDs(contact models.Contact) []string {
	ids := make([]string, 0, len(contact.CampaignIDs)+len(contact.CampaignMemberships)+1)
	seen := map[string]bool{}

	if trimmed := strings.TrimSpace(contact.CampaignID); trimmed != "" {
		ids = append(ids, trimmed)
		seen[trimmed] = true
	}

	for _, campaignID := range contact.CampaignIDs {
		trimmed := strings.TrimSpace(campaignID)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		ids = append(ids, trimmed)
		seen[trimmed] = true
	}

	for _, membership := range contact.CampaignMemberships {
		trimmed := strings.TrimSpace(membership.CampaignID)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		ids = append(ids, trimmed)
		seen[trimmed] = true
	}

	return ids
}

func normalizeCampaignMembershipStatus(status models.CampaignContactStatus) models.CampaignContactStatus {
	switch strings.ToLower(string(status)) {
	case string(models.CampaignContactStatusQueued):
		return models.CampaignContactStatusQueued
	case string(models.CampaignContactStatusInProgress):
		return models.CampaignContactStatusInProgress
	case string(models.CampaignContactStatusSuccess):
		return models.CampaignContactStatusSuccess
	case string(models.CampaignContactStatusFailed):
		return models.CampaignContactStatusFailed
	default:
		return models.CampaignContactStatusQueued
	}
}

func normalizeCampaignMemberships(contact models.Contact) []models.ContactCampaignMembership {
	now := contact.UpdatedAt
	if strings.TrimSpace(now) == "" {
		now = time.Now().UTC().Format(time.RFC3339)
	}

	memberships := make([]models.ContactCampaignMembership, 0, len(contact.CampaignMemberships)+len(contact.CampaignIDs)+1)
	seen := map[string]bool{}

	for _, membership := range contact.CampaignMemberships {
		trimmed := strings.TrimSpace(membership.CampaignID)
		if trimmed == "" || seen[trimmed] {
			continue
		}

		normalizedMembership := membership
		normalizedMembership.CampaignID = trimmed
		normalizedMembership.Status = normalizeCampaignMembershipStatus(normalizedMembership.Status)
		if strings.TrimSpace(normalizedMembership.CreatedAt) == "" {
			normalizedMembership.CreatedAt = now
		}
		if strings.TrimSpace(normalizedMembership.UpdatedAt) == "" {
			normalizedMembership.UpdatedAt = normalizedMembership.CreatedAt
		}

		memberships = append(memberships, normalizedMembership)
		seen[trimmed] = true
	}

	legacyCampaignIDs := normalizeCampaignIDs(models.Contact{
		CampaignID:  contact.CampaignID,
		CampaignIDs: contact.CampaignIDs,
	})
	for _, campaignID := range legacyCampaignIDs {
		if seen[campaignID] {
			continue
		}

		memberships = append(memberships, models.ContactCampaignMembership{
			CampaignID: campaignID,
			Status:     models.CampaignContactStatusQueued,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		seen[campaignID] = true
	}

	return memberships
}
