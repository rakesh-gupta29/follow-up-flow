import { TableCell } from "@/components/ui/table"
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
    cell: (campaign) => (
      <TableCell key="name">
        <span className="font-medium">{campaign.name}</span>
      </TableCell>
    ),
  },
  {
    id: "description",
    header: "Description",
    cell: (campaign) => (
      <TableCell key="description">
        <span className="text-muted-foreground">{campaign.description ?? "—"}</span>
      </TableCell>
    ),
  },
  {
    id: "created_at",
    header: "Created",
    cell: (campaign) => (
      <TableCell key="created_at">
        {new Date(campaign.created_at).toLocaleDateString()}
      </TableCell>
    ),
  },
  {
    id: "contacts_count",
    header: "Contacts",
    cell: (campaign) => (
      <TableCell key="contacts_count">
        {campaign.contacts_count ?? 0}
      </TableCell>
    ),
  },
]
