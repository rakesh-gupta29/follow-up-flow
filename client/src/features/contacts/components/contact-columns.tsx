import type { ReactNode } from "react"

import type { ContactListItem } from "../../../types/contact"

export type ContactColumn<TData> = {
  id: string
  header: string
  className?: string
  cell: (row: TData) => ReactNode
}

export const contactColumns: ContactColumn<ContactListItem>[] = [
  {
    id: "name",
    header: "Person",
    cell: (contact) => `${contact.first_name} ${contact.last_name}`.trim() || contact.email,
  },
  {
    id: "company",
    header: "Company",
    cell: (contact) => contact.company || "No company",
  },
  {
    id: "campaigns",
    header: "Campaigns",
    cell: (contact) =>
      contact.campaigns && contact.campaigns.length > 0
        ? `${contact.campaigns.length} campaign${contact.campaigns.length === 1 ? "" : "s"}`
        : "Unassigned",
  },
  {
    id: "status",
    header: "Status",
    cell: (contact) => contact.status,
  },
]

export const campaignContactColumns: ContactColumn<ContactListItem>[] = [
  {
    id: "name",
    header: "Person",
    cell: (contact) => `${contact.first_name} ${contact.last_name}`.trim() || contact.email,
  },
  {
    id: "company",
    header: "Company",
    cell: (contact) => contact.company || "No company",
  },
  {
    id: "current-campaign",
    header: "Current campaign",
    cell: (contact) => contact.campaign?.name || "Unassigned",
  },
  {
    id: "current-state",
    header: "Current state",
    cell: (contact) => contact.campaign?.status || "Unknown",
  },
  {
    id: "logs",
    header: "Logs",
    cell: () => "View logs",
  },
]
