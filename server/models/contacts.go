package models

type ContactStatus string

const (
	ContactStatusActive       ContactStatus = "active"
	ContactStatusUnsubscribed ContactStatus = "unsubscribed"
	ContactStatusBounced      ContactStatus = "bounced"
)

type ContactCampaign struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type Contact struct {
	ID          string            `json:"id" bson:"id"`
	Email       string            `json:"email" bson:"email"`
	FirstName   string            `json:"first_name" bson:"first_name"`
	LastName    string            `json:"last_name" bson:"last_name"`
	Phone       string            `json:"phone,omitempty" bson:"phone,omitempty"`
	Company     string            `json:"company,omitempty" bson:"company,omitempty"`
	Tags        []string          `json:"tags,omitempty" bson:"tags,omitempty"`
	CampaignID  string            `json:"campaign_id,omitempty" bson:"campaign_id,omitempty"`
	CampaignIDs []string          `json:"campaign_ids,omitempty" bson:"campaign_ids,omitempty"`
	Status      ContactStatus     `json:"status" bson:"status"`
	Meta        map[string]string `json:"meta,omitempty" bson:"meta,omitempty"`
	CreatedAt   string            `json:"created_at" bson:"created_at"`
	UpdatedAt   string            `json:"updated_at" bson:"updated_at"`
}

type ContactListItem struct {
	Contact   `bson:",inline"`
	Campaign  *ContactCampaign  `json:"campaign,omitempty" bson:"campaign,omitempty"`
	Campaigns []ContactCampaign `json:"campaigns,omitempty" bson:"campaigns,omitempty"`
}
