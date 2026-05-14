package models

type ContactStatus string

const (
	ContactStatusActive       ContactStatus = "active"
	ContactStatusUnsubscribed ContactStatus = "unsubscribed"
	ContactStatusBounced      ContactStatus = "bounced"
	ContactStatusCallback     ContactStatus = "callback"
)

type CampaignContactStatus string

const (
	CampaignContactStatusQueued     CampaignContactStatus = "queued"
	CampaignContactStatusInProgress CampaignContactStatus = "in_progress"
	CampaignContactStatusSuccess    CampaignContactStatus = "success"
	CampaignContactStatusFailed     CampaignContactStatus = "failed"
	CampaignContactStatusCallback   CampaignContactStatus = "callback"
)

type ContactCampaign struct {
	ID     string                `json:"id" bson:"id"`
	Name   string                `json:"name" bson:"name"`
	Status CampaignContactStatus `json:"status" bson:"status"`
}

type ContactCampaignMembership struct {
	CampaignID string                `json:"campaign_id" bson:"campaign_id"`
	Status     CampaignContactStatus `json:"status" bson:"status"`
	CreatedAt  string                `json:"created_at" bson:"created_at"`
	UpdatedAt  string                `json:"updated_at" bson:"updated_at"`
	Campaign   *Campaign             `json:"campaign,omitempty" bson:"-"`
}

type ContactCampaignLog struct {
	CampaignID string                `json:"campaign_id" bson:"campaign_id"`
	Status     CampaignContactStatus `json:"status" bson:"status"`
	Action     string                `json:"action" bson:"action"`
	CreatedAt  string                `json:"created_at" bson:"created_at"`
}

type Contact struct {
	ID                  string                      `json:"id" bson:"id"`
	Email               string                      `json:"email" bson:"email"`
	FirstName           string                      `json:"first_name" bson:"first_name"`
	LastName            string                      `json:"last_name" bson:"last_name"`
	Phone               string                      `json:"phone,omitempty" bson:"phone,omitempty"`
	PropertyName        string                      `json:"property_name,omitempty" bson:"property_name,omitempty"`
	QuestionnaireURL    string                      `json:"questionnaire_url,omitempty" bson:"questionnaire_url,omitempty"`
	ThreadID            string                      `json:"thread_id,omitempty" bson:"thread_id,omitempty"`
	CallID              string                      `json:"call_id,omitempty" bson:"call_id,omitempty"`
	CampaignID          string                      `json:"campaign_id,omitempty" bson:"campaign_id,omitempty"`
	CampaignIDs         []string                    `json:"campaign_ids,omitempty" bson:"campaign_ids,omitempty"`
	CampaignMemberships []ContactCampaignMembership `json:"campaign_memberships,omitempty" bson:"campaign_memberships,omitempty"`
	CampaignLogs        []ContactCampaignLog        `json:"campaign_logs,omitempty" bson:"campaign_logs,omitempty"`
	Status              ContactStatus               `json:"status" bson:"status"`
	CreatedAt           string                      `json:"created_at" bson:"created_at"`
	UpdatedAt           string                      `json:"updated_at" bson:"updated_at"`
	DeletedAt           string                      `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type ContactListItem struct {
	Contact   `bson:",inline"`
	Campaign  *ContactCampaign  `json:"campaign,omitempty" bson:"campaign,omitempty"`
	Campaigns []ContactCampaign `json:"campaigns,omitempty" bson:"campaigns,omitempty"`
}
