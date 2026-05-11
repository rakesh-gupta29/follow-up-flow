package repository

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/shingo/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CampaignsRepository struct {
	client              *mongo.Client
	campaignsCollection *mongo.Collection
	contactsCollection  *mongo.Collection
}

func NewCampaignsRepository(db *mongo.Client) *CampaignsRepository {
	return &CampaignsRepository{
		client:              db,
		campaignsCollection: db.Database(appDatabaseName).Collection("campaigns"),
		contactsCollection:  db.Database(appDatabaseName).Collection("contacts"),
	}
}

func (r *CampaignsRepository) EnsureCollection(ctx context.Context) error {
	db := r.client.Database(appDatabaseName)

	collections, err := db.ListCollectionNames(ctx, bson.M{"name": "campaigns"})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		if err := db.CreateCollection(ctx, "campaigns"); err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "exists") {
				return err
			}
		}
	}

	_, err = r.campaignsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
	})
	if err != nil {
		return err
	}

	_, err = r.contactsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "campaign_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "campaign_ids", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "campaign_memberships.campaign_id", Value: 1}},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *CampaignsRepository) ListCampaigns(ctx context.Context) ([]models.CampaignListItem, error) {
	cursor, err := r.campaignsCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{
		{Key: "stage", Value: 1},
		{Key: "created_at", Value: -1},
	}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	campaigns := make([]models.CampaignListItem, 0)
	for cursor.Next(ctx) {
		var campaign models.Campaign
		if err := cursor.Decode(&campaign); err != nil {
			return nil, err
		}

		contactsCount, err := r.contactsCollection.CountDocuments(ctx, bson.M{
			"$or": []bson.M{
				{"campaign_id": campaign.ID},
				{"campaign_ids": campaign.ID},
				{"campaign_memberships.campaign_id": campaign.ID},
			},
		})
		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, models.CampaignListItem{
			Campaign:      campaign,
			ContactsCount: contactsCount,
		})
	}

	return campaigns, cursor.Err()
}

func (r *CampaignsRepository) GetCampaign(ctx context.Context, id string) (*models.Campaign, error) {
	var campaign models.Campaign
	err := r.campaignsCollection.FindOne(ctx, bson.M{"id": id}).Decode(&campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (r *CampaignsRepository) ListCampaignContacts(ctx context.Context, campaignID string, page, limit int64, search string, status string) ([]models.ContactListItem, int64, error) {
	campaign, err := r.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"campaign_id": campaignID},
					{"campaign_ids": campaignID},
					{"campaign_memberships.campaign_id": campaignID},
				},
			},
		},
	}

	if trimmedSearch := strings.TrimSpace(search); trimmedSearch != "" {
		pattern := regexp.QuoteMeta(trimmedSearch)
		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"$or": []bson.M{
				{"email": bson.M{"$regex": pattern, "$options": "i"}},
				{"first_name": bson.M{"$regex": pattern, "$options": "i"}},
				{"last_name": bson.M{"$regex": pattern, "$options": "i"}},
				{"property_name": bson.M{"$regex": pattern, "$options": "i"}},
				{"phone": bson.M{"$regex": pattern, "$options": "i"}},
				{"questionnaire_url": bson.M{"$regex": pattern, "$options": "i"}},
				{"thread_id": bson.M{"$regex": pattern, "$options": "i"}},
			},
		})
	}

	if normalizedStatus := normalizeStatus(models.ContactStatus(status)); normalizedStatus != "" {
		filter["status"] = normalizedStatus
	}

	total, err := r.contactsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.contactsCollection.Find(ctx, filter, options.Find().
		SetSkip((page-1)*limit).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}}))
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

		if err := r.enrichContactCampaignMemberships(ctx, &contact); err != nil {
			return nil, 0, err
		}

		contacts = append(contacts, models.ContactListItem{
			Contact: contact,
			Campaign: &models.ContactCampaign{
				ID:     campaign.ID,
				Name:   campaign.Name,
				Status: findCampaignStatus(contact, campaign.ID),
			},
			Campaigns: []models.ContactCampaign{
				{
					ID:     campaign.ID,
					Name:   campaign.Name,
					Status: findCampaignStatus(contact, campaign.ID),
				},
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (r *CampaignsRepository) AttachContact(ctx context.Context, campaignID, contactID string) (*models.Contact, error) {
	if _, err := r.GetCampaign(ctx, campaignID); err != nil {
		return nil, err
	}

	contact, err := r.getContact(ctx, contactID)
	if err != nil {
		return nil, err
	}

	updatedContact, _, err := r.ensureCampaignQueued(ctx, contact, campaignID)
	if err == nil && updatedContact != nil {
		if enrichErr := r.enrichContactCampaignMemberships(ctx, updatedContact); enrichErr != nil {
			return nil, enrichErr
		}
	}
	return updatedContact, err
}

func (r *CampaignsRepository) AttachContacts(ctx context.Context, campaignID string, contactIDs []string) (int64, error) {
	if _, err := r.GetCampaign(ctx, campaignID); err != nil {
		return 0, err
	}

	if len(contactIDs) == 0 {
		return 0, nil
	}

	trimmedIDs := make([]string, 0, len(contactIDs))
	for _, contactID := range contactIDs {
		if trimmed := strings.TrimSpace(contactID); trimmed != "" {
			trimmedIDs = append(trimmedIDs, trimmed)
		}
	}

	if len(trimmedIDs) == 0 {
		return 0, nil
	}

	var modifiedCount int64
	for _, contactID := range trimmedIDs {
		contact, err := r.getContact(ctx, contactID)
		if err != nil {
			return modifiedCount, err
		}

		_, attached, err := r.ensureCampaignQueued(ctx, contact, campaignID)
		if err != nil {
			return modifiedCount, err
		}
		if attached {
			modifiedCount++
		}
	}

	return modifiedCount, nil
}

func (r *CampaignsRepository) RemoveContact(ctx context.Context, campaignID, contactID string) (*models.Contact, error) {
	if _, err := r.GetCampaign(ctx, campaignID); err != nil {
		return nil, err
	}

	contact, err := r.getContact(ctx, contactID)
	if err != nil {
		return nil, err
	}

	memberships := normalizeCampaignMemberships(*contact)
	index := findCampaignMembershipIndex(memberships, campaignID)
	if index == -1 {
		return nil, mongo.ErrNoDocuments
	}

	status := memberships[index].Status
	memberships = append(memberships[:index], memberships[index+1:]...)

	contact.CampaignMemberships = memberships
	contact.CampaignIDs = campaignIDsFromMemberships(memberships)
	if len(contact.CampaignIDs) > 0 {
		contact.CampaignID = contact.CampaignIDs[0]
	} else {
		contact.CampaignID = ""
	}
	contact.CampaignLogs = append(contact.CampaignLogs, models.ContactCampaignLog{
		CampaignID: campaignID,
		Status:     status,
		Action:     "removed",
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	})

	if err := r.saveCampaignState(ctx, contact); err != nil {
		return nil, err
	}

	if err := r.enrichContactCampaignMemberships(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *CampaignsRepository) UpdateCampaignContactStatus(ctx context.Context, campaignID, contactID string, status models.CampaignContactStatus, nextCampaignID string) (*models.Contact, error) {
	if _, err := r.GetCampaign(ctx, campaignID); err != nil {
		return nil, err
	}

	if trimmedNextCampaignID := strings.TrimSpace(nextCampaignID); trimmedNextCampaignID != "" {
		if _, err := r.GetCampaign(ctx, trimmedNextCampaignID); err != nil {
			return nil, err
		}
		nextCampaignID = trimmedNextCampaignID
	}

	contact, err := r.getContact(ctx, contactID)
	if err != nil {
		return nil, err
	}

	memberships := normalizeCampaignMemberships(*contact)
	index := findCampaignMembershipIndex(memberships, campaignID)
	if index == -1 {
		return nil, mongo.ErrNoDocuments
	}

	now := time.Now().UTC().Format(time.RFC3339)
	memberships[index].Status = normalizeCampaignMembershipStatus(status)
	memberships[index].UpdatedAt = now
	contact.CampaignLogs = append(contact.CampaignLogs, models.ContactCampaignLog{
		CampaignID: campaignID,
		Status:     memberships[index].Status,
		Action:     "status_updated",
		CreatedAt:  now,
	})

	if nextCampaignID != "" {
		nextIndex := findCampaignMembershipIndex(memberships, nextCampaignID)
		if nextCampaignID != campaignID && nextIndex == -1 {
			memberships = append(memberships, models.ContactCampaignMembership{
				CampaignID: nextCampaignID,
				Status:     models.CampaignContactStatusQueued,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
			contact.CampaignLogs = append(contact.CampaignLogs, models.ContactCampaignLog{
				CampaignID: nextCampaignID,
				Status:     models.CampaignContactStatusQueued,
				Action:     "queued_next_campaign",
				CreatedAt:  now,
			})
		}
	}

	contact.CampaignMemberships = memberships
	contact.CampaignIDs = campaignIDsFromMemberships(memberships)
	if len(contact.CampaignIDs) > 0 {
		contact.CampaignID = contact.CampaignIDs[0]
	}
	contact.UpdatedAt = now

	if err := r.saveCampaignState(ctx, contact); err != nil {
		return nil, err
	}

	if err := r.enrichContactCampaignMemberships(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *CampaignsRepository) GetQueuedCampaignContacts(ctx context.Context, campaignID string) ([]models.ContactListItem, error) {
	campaign, err := r.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"$or": []bson.M{
			bson.M{
				"campaign_memberships": bson.M{
					"$elemMatch": bson.M{
						"campaign_id": campaignID,
						"status":      models.CampaignContactStatusQueued,
					},
				},
			},
			bson.M{
				"$and": []bson.M{
					{
						"$or": []bson.M{
							{"campaign_id": campaignID},
							{"campaign_ids": campaignID},
						},
					},
					{"campaign_memberships": bson.M{"$exists": false}},
				},
			},
		},
	}

	cursor, err := r.contactsCollection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	contacts := make([]models.ContactListItem, 0)
	for cursor.Next(ctx) {
		var contact models.Contact
		if err := cursor.Decode(&contact); err != nil {
			return nil, err
		}

		if err := r.enrichContactCampaignMemberships(ctx, &contact); err != nil {
			return nil, err
		}

		status := findCampaignStatus(contact, campaignID)
		if status != models.CampaignContactStatusQueued {
			continue
		}

		contacts = append(contacts, models.ContactListItem{
			Contact: contact,
			Campaign: &models.ContactCampaign{
				ID:     campaign.ID,
				Name:   campaign.Name,
				Status: status,
			},
			Campaigns: []models.ContactCampaign{
				{
					ID:     campaign.ID,
					Name:   campaign.Name,
					Status: status,
				},
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *CampaignsRepository) getContact(ctx context.Context, contactID string) (*models.Contact, error) {
	var contact models.Contact
	err := r.contactsCollection.FindOne(ctx, bson.M{"id": contactID}).Decode(&contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *CampaignsRepository) ensureCampaignQueued(ctx context.Context, contact *models.Contact, campaignID string) (*models.Contact, bool, error) {
	memberships := normalizeCampaignMemberships(*contact)
	now := time.Now().UTC().Format(time.RFC3339)
	action := "attached"
	status := models.CampaignContactStatusQueued

	if index := findCampaignMembershipIndex(memberships, campaignID); index != -1 {
		memberships[index].Status = status
		if strings.TrimSpace(memberships[index].CreatedAt) == "" {
			memberships[index].CreatedAt = now
		}
		memberships[index].UpdatedAt = now
		action = "requeued"
	} else {
		memberships = append(memberships, models.ContactCampaignMembership{
			CampaignID: campaignID,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	contact.CampaignMemberships = memberships
	contact.CampaignIDs = campaignIDsFromMemberships(memberships)
	if len(contact.CampaignIDs) > 0 {
		contact.CampaignID = contact.CampaignIDs[0]
	}
	contact.CampaignLogs = append(contact.CampaignLogs, models.ContactCampaignLog{
		CampaignID: campaignID,
		Status:     status,
		Action:     action,
		CreatedAt:  now,
	})
	contact.UpdatedAt = now

	if err := r.saveCampaignState(ctx, contact); err != nil {
		return nil, false, err
	}

	if err := r.enrichContactCampaignMemberships(ctx, contact); err != nil {
		return nil, false, err
	}

	return contact, true, nil
}

func (r *CampaignsRepository) saveCampaignState(ctx context.Context, contact *models.Contact) error {
	_, err := r.contactsCollection.UpdateOne(ctx, bson.M{"id": contact.ID}, bson.M{
		"$set": bson.M{
			"campaign_id":          contact.CampaignID,
			"campaign_ids":         contact.CampaignIDs,
			"campaign_memberships": contact.CampaignMemberships,
			"campaign_logs":        contact.CampaignLogs,
			"updated_at":           contact.UpdatedAt,
		},
	})
	return err
}

func findCampaignMembershipIndex(memberships []models.ContactCampaignMembership, campaignID string) int {
	for index, membership := range memberships {
		if membership.CampaignID == campaignID {
			return index
		}
	}

	return -1
}

func campaignIDsFromMemberships(memberships []models.ContactCampaignMembership) []string {
	ids := make([]string, 0, len(memberships))
	for _, membership := range memberships {
		if trimmed := strings.TrimSpace(membership.CampaignID); trimmed != "" {
			ids = append(ids, trimmed)
		}
	}

	return ids
}

func findCampaignStatus(contact models.Contact, campaignID string) models.CampaignContactStatus {
	memberships := normalizeCampaignMemberships(contact)
	for _, membership := range memberships {
		if membership.CampaignID == campaignID {
			return membership.Status
		}
	}

	return models.CampaignContactStatusQueued
}

func (r *CampaignsRepository) enrichContactCampaignMemberships(ctx context.Context, contact *models.Contact) error {
	memberships := normalizeCampaignMemberships(*contact)
	for index := range memberships {
		var campaignDoc models.Campaign
		err := r.campaignsCollection.FindOne(ctx, bson.M{"id": memberships[index].CampaignID}).Decode(&campaignDoc)
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}
		if err == nil {
			memberships[index].Campaign = &campaignDoc
		}
	}

	contact.CampaignMemberships = memberships
	return nil
}
