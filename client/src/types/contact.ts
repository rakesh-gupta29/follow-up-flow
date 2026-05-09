export type ContactStatus = "active" | "unsubscribed" | "bounced"

export type Contact = {
  id: string
  email: string
  first_name: string
  last_name: string
  phone?: string
  company?: string
  tags?: string[]
  status: ContactStatus
  meta?: Record<string, string>
  created_at: string
  updated_at: string
}

export type ContactCampaign = {
  id: string
  name: string
}

export type ContactListItem = Contact & {
  campaigns?: ContactCampaign[]
}
