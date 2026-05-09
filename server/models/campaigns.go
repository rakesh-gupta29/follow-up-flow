package models

type Campaign struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   string `json:"created_at" bson:"created_at"`
	UpdatedAt   string `json:"updated_at" bson:"updated_at"`
}

type CampaignListItem struct {
	Campaign
	ContactsCount int64 `json:"contacts_count"`
}
