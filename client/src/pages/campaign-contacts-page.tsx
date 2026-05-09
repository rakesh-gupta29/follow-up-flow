import { useMemo, useState } from "react"
import { useParams } from "react-router-dom"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useContactsQuery } from "../features/contacts/api/get-contacts"
import { ContactsTable } from "../features/contacts/components/contacts-table"

export function CampaignContactsPage() {
  const { campaignId = "" } = useParams()
  const [selectedContactIds, setSelectedContactIds] = useState<string[]>([])
  const { data, isLoading } = useContactsQuery({
    page: 1,
    limit: 25,
    search: "",
    campaignId,
  })

  const contacts = data?.items ?? []

  const selectedOnPage = useMemo(
    () => contacts.filter((contact) => selectedContactIds.includes(contact.id)).length,
    [contacts, selectedContactIds]
  )

  const handleToggleContact = (contactId: string, checked: boolean) => {
    setSelectedContactIds((current) =>
      checked ? [...new Set([...current, contactId])] : current.filter((id) => id !== contactId)
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Campaign contacts</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <ContactsTable
          contacts={contacts}
          selectedContactIds={selectedContactIds}
          onToggleContact={handleToggleContact}
        />
        <p className="text-muted-foreground text-sm">
          {isLoading ? "Loading contacts..." : `${selectedOnPage} selected on this page`}
        </p>
      </CardContent>
    </Card>
  )
}
