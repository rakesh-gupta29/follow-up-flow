export type ContactStatus = "active" | "unsubscribed" | "bounced"
export type CampaignContactStatus = "queued" | "in_progress" | "success" | "failed"

export type Contact = {
  id: string
  email: string
  first_name: string
  last_name: string
  phone?: string
  property_name?: string
  questionnaire_url?: string
  thread_id?: string
  campaign_id?: string
  campaign_ids?: string[]
  status: ContactStatus
  created_at: string
  updated_at: string
  campaign_memberships?: ContactCampaignMembership[]
  campaign_logs?: ContactCampaignLog[]
}

export type ContactCampaign = {
  id: string
  name: string
  status: CampaignContactStatus
}

export type ContactCampaignLog = {
  campaign_id: string
  status: CampaignContactStatus
  action: string
  created_at: string
}

export type ContactCampaignMembership = {
  campaign_id: string
  status: CampaignContactStatus
  created_at: string
  updated_at: string
  name?: string
  campaign?: {
    id?: string
    name?: string
    stage?: number
  }
}

export type ContactListItem = Contact & {
  campaign?: ContactCampaign
  campaigns?: ContactCampaign[]
}
