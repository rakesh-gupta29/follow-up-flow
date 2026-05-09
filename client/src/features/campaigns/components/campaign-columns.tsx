import type { ReactNode } from "react"

import type { Campaign } from "../../../types/campaign"

export type CampaignColumn<TData> = {
  id: string
  header: string
  className?: string
  cell: (row: TData) => ReactNode
}

export const campaignColumns: CampaignColumn<Campaign>[] = [
  {
    id: "name",
    header: "Campaign",
    cell: (campaign) => campaign.name,
  },
  {
    id: "status",
    header: "Status",
    cell: (campaign) => campaign.status,
  },
  {
    id: "stages",
    header: "Stages",
    cell: (campaign) => campaign.stages.length,
  },
  {
    id: "updated_at",
    header: "Updated",
    cell: (campaign) => new Date(campaign.updated_at).toLocaleDateString(),
  },
]
