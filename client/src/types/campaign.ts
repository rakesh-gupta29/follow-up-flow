export type CampaignStatus = "draft" | "active" | "paused" | "archived"
export type StageType = "email" | "wait" | "webhook"

export type Stage = {
  stage_number: number
  name: string
  type: StageType
  delay_hours: number
  n8n_webhook_url: string
  template_id?: string
  max_retries: number
}

export type Campaign = {
  id: string
  name: string
  description?: string
  status: CampaignStatus
  stages: Stage[]
  created_at: string
  updated_at: string
}
