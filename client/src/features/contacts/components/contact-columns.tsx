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
    id: "campaign",
    header: "Campaign",
    cell: (contact) => contact.campaign?.name || "Unassigned",
  },
  {
    id: "status",
    header: "Status",
    cell: (contact) => contact.status,
  },
]
