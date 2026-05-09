import { useMemo, useState } from "react"
import { useParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogBody,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { useContactsQuery } from "../features/contacts/api/get-contacts"
import { useRemoveContactsFromCampaignMutation } from "../features/contacts/api/remove-contacts-from-campaign"
import { ContactsTable } from "../features/contacts/components/contacts-table"

export function CampaignContactsPage() {
  const { campaignId = "" } = useParams()
  const [selectedContactIds, setSelectedContactIds] = useState<string[]>([])
  const [isConfirmOpen, setIsConfirmOpen] = useState(false)
  const removeContactsMutation = useRemoveContactsFromCampaignMutation()
  const { data, isLoading, isError, error, isFetching } = useContactsQuery({
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

  const handleRemoveContacts = async () => {
    if (!campaignId || selectedContactIds.length === 0) {
      return
    }

    await removeContactsMutation.mutateAsync({
      campaignId,
      contactIds: selectedContactIds,
    })
    setSelectedContactIds([])
    setIsConfirmOpen(false)
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between gap-10">
          <CardTitle>Campaign contacts</CardTitle>
          <div className="flex  items-center gap-3">
            <Button
              type="button"
              variant="outline"
              onClick={() => setIsConfirmOpen(true)}
              disabled={selectedContactIds.length === 0 || removeContactsMutation.isPending}
            >
              {removeContactsMutation.isPending ? "Removing..." : "Remove from campaign"}
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {isLoading ? (
          <div className="text-muted-foreground rounded-xl border border-dashed p-8 text-sm">
            Loading contacts...
          </div>
        ) : isError ? (
          <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-4 text-sm text-destructive">
            {error instanceof Error
              ? error.message
              : "Failed to load campaign contacts."}
          </div>
        ) : contacts.length === 0 ? (
          <div className="text-muted-foreground rounded-xl border border-dashed p-8 text-sm">
            No contacts found for this campaign.
          </div>
        ) : (
          <ContactsTable
            contacts={contacts}
            selectedContactIds={selectedContactIds}
            onToggleContact={handleToggleContact}
            variant="campaign"
          />
        )}
        {removeContactsMutation.isError ? (
          <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-4 text-sm text-destructive">
            {removeContactsMutation.error instanceof Error
              ? removeContactsMutation.error.message
              : "Failed to remove contacts from the campaign."}
          </div>
        ) : null}
        <p className="text-muted-foreground text-sm">
          {isLoading
            ? "Loading contacts..."
            : isFetching
              ? "Refreshing contacts..."
              : `${selectedOnPage} selected on this page`}
        </p>
      </CardContent>

      <Dialog
        open={isConfirmOpen}
        onOpenChange={(open) => {
          if (!removeContactsMutation.isPending) {
            setIsConfirmOpen(open)
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Remove contacts from campaign</DialogTitle>
            <DialogDescription>
              This will remove {selectedContactIds.length} selected contact
              {selectedContactIds.length === 1 ? "" : "s"} from this campaign.
            </DialogDescription>
          </DialogHeader>
          <DialogBody className="text-muted-foreground text-sm">
            This action updates the campaign membership immediately.
          </DialogBody>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setIsConfirmOpen(false)}
              disabled={removeContactsMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              type="button"
              onClick={handleRemoveContacts}
              disabled={removeContactsMutation.isPending}
            >
              {removeContactsMutation.isPending ? "Removing..." : "Confirm removal"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  )
}
