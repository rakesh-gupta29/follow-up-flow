import { useMemo, useState } from "react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Select } from "@/components/ui/select"
import { useContactsQuery } from "../features/contacts/api/get-contacts"
import { AddToCampaignDialog } from "../features/contacts/components/add-to-campaign-dialog"
import { ContactsTable } from "../features/contacts/components/contacts-table"
import type { ContactStatus } from "../types/contact"

const limits = [25, 50, 100]

export function ContactsPage() {
  const [page, setPage] = useState(1)
  const [limit, setLimit] = useState(25)
  const [search, setSearch] = useState("")
  const [status, setStatus] = useState<ContactStatus | undefined>(undefined)
  const [selectedContactIds, setSelectedContactIds] = useState<string[]>([])
  const [dialogOpen, setDialogOpen] = useState(false)
  const { data, isLoading } = useContactsQuery({
    page,
    limit,
    search,
    status,
  })

  const contacts = data?.items ?? []
  const currentPage = data?.pagination?.page ?? 1
  const totalItems = data?.pagination.total ?? 1

  const allVisibleSelected = useMemo(
    () => contacts.length > 0 && contacts.every((contact) => selectedContactIds.includes(contact.id)),
    [contacts, selectedContactIds]
  )

  const handleToggleContact = (contactId: string, checked: boolean) => {
    setSelectedContactIds((current) =>
      checked ? [...new Set([...current, contactId])] : current.filter((id) => id !== contactId)
    )
  }

  const handleSelectAll = () => {
    if (allVisibleSelected) {
      setSelectedContactIds((current) =>
        current.filter((id) => !contacts.some((contact) => contact.id === id))
      )
      return
    }

    setSelectedContactIds((current) => [
      ...new Set([...current, ...contacts.map((contact) => contact.id)]),
    ])
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader className="flex">
          <div className="flex justify-between  w-full ">
            <div>
              <CardTitle>Contacts</CardTitle>
              <p className="text-muted-foreground text-sm">
                {isLoading ? "Loading contacts..." : `${totalItems} contacts`}
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <Button type="button" variant="outline" onClick={handleSelectAll}>
                {allVisibleSelected ? "Clear selection" : "Select all"}
              </Button>
              <Button
                type="button"
                onClick={() => setDialogOpen(true)}
                disabled={selectedContactIds.length === 0}
              >
                Add to campaign
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 ">


          <div className="flex flex-wrap items-center gap-3">
            <Input
              className="w-full max-w-sm"
              placeholder="Search contacts"
              value={search}
              onChange={(event) => {
                setPage(1)
                setSearch(event.target.value)
              }}
            />
            <Select
              value={status ?? ""}
              onChange={(event) => {
                setPage(1)
                setStatus((event.target.value || undefined) as ContactStatus | undefined)
              }}
            >
              <option value="">All statuses</option>
              <option value="active">Active</option>
              <option value="unsubscribed">Unsubscribed</option>
              <option value="bounced">Bounced</option>
            </Select>

          </div>

          <ContactsTable
            contacts={contacts}
            selectedContactIds={selectedContactIds}
            onToggleContact={handleToggleContact}
          />

          <div className="flex items-center justify-end gap-6">
            <Select
              value={String(limit)}
              onChange={(event) => {
                setPage(1)
                setLimit(Number(event.target.value))
              }}
            >
              {limits.map((value) => (
                <option key={value} value={value}>
                  {value} per page
                </option>
              ))}
            </Select>
            <div className="flex items-center gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setPage((current) => Math.max(current - 1, 1))}
                disabled={page <= 1}
              >
                Previous
              </Button>
              <span className="text-sm">Page {page}</span>
              <Button
                type="button"
                variant="outline"
                onClick={() => setPage((current) => Math.min(current + 1, totalItems))}
                disabled={page >= currentPage}
              >
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <AddToCampaignDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        selectedContactIds={selectedContactIds}
      />
    </div>
  )
}
