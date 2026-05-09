package repository

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
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
	})
	if err != nil {
		return err
	}

	return r.seedDefaultCampaigns(ctx)
}

func (r *CampaignsRepository) seedDefaultCampaigns(ctx context.Context) error {
	count, err := r.campaignsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	campaigns := []any{
		models.Campaign{
			ID:          uuid.NewString(),
			Name:        "Welcome Campaign",
			Description: "Default onboarding campaign",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		models.Campaign{
			ID:          uuid.NewString(),
			Name:        "Re-Engagement Campaign",
			Description: "Default re-engagement campaign",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		models.Campaign{
			ID:          uuid.NewString(),
			Name:        "Product Update Campaign",
			Description: "Default product update campaign",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	_, err = r.campaignsCollection.InsertMany(ctx, campaigns)
	return err
}

func (r *CampaignsRepository) ListCampaigns(ctx context.Context) ([]models.CampaignListItem, error) {
	cursor, err := r.campaignsCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
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
				{"company": bson.M{"$regex": pattern, "$options": "i"}},
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

		contacts = append(contacts, models.ContactListItem{
			Contact: contact,
			Campaign: &models.ContactCampaign{
				ID:   campaign.ID,
				Name: campaign.Name,
			},
			Campaigns: []models.ContactCampaign{
				{
					ID:   campaign.ID,
					Name: campaign.Name,
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

	update := bson.M{
		"$addToSet": bson.M{
			"campaign_ids": campaignID,
		},
		"$set": bson.M{
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		},
	}

	result := r.contactsCollection.FindOneAndUpdate(
		ctx,
		bson.M{"id": contactID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var contact models.Contact
	if err := result.Decode(&contact); err != nil {
		return nil, err
	}

	return &contact, nil
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

	result, err := r.contactsCollection.UpdateMany(
		ctx,
		bson.M{"id": bson.M{"$in": trimmedIDs}},
		bson.M{
			"$addToSet": bson.M{
				"campaign_ids": bson.M{"$each": trimmedIDsToCampaigns(campaignID)},
			},
			"$set": bson.M{
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			},
		},
	)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (r *CampaignsRepository) RemoveContact(ctx context.Context, campaignID, contactID string) (*models.Contact, error) {
	if _, err := r.GetCampaign(ctx, campaignID); err != nil {
		return nil, err
	}

	result := r.contactsCollection.FindOneAndUpdate(
		ctx,
		bson.M{
			"id": contactID,
			"$or": []bson.M{
				{"campaign_id": campaignID},
				{"campaign_ids": campaignID},
			},
		},
		bson.M{
			"$set": bson.M{
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			},
			"$pull": bson.M{
				"campaign_ids": campaignID,
			},
			"$unset": bson.M{
				"campaign_id": "",
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var contact models.Contact
	if err := result.Decode(&contact); err != nil {
		return nil, err
	}

	return &contact, nil
}

func trimmedIDsToCampaigns(campaignID string) []string {
	if strings.TrimSpace(campaignID) == "" {
		return []string{}
	}

	return []string{strings.TrimSpace(campaignID)}
}
