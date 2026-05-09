import { Badge } from "@/components/ui/badge"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Dialog,
  DialogBody,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useState } from "react"
import { contactColumns } from "./contact-columns"
import type { ContactListItem } from "../../../types/contact"

type ContactsTableProps = {
  contacts: ContactListItem[]
  selectedContactIds: string[]
  onToggleContact: (contactId: string, checked: boolean) => void
  onRowClick?: (contact: ContactListItem) => void
}

function getStatusVariant(status: ContactListItem["status"]) {
  if (status === "active") {
    return "default"
  }

  if (status === "unsubscribed") {
    return "secondary"
  }

  return "outline"
}

export function ContactsTable({
  contacts,
  selectedContactIds,
  onToggleContact,
  onRowClick,
}: ContactsTableProps) {
  const [campaignDetailsContact, setCampaignDetailsContact] = useState<ContactListItem | null>(null)

  return (
    <>
      <div className="overflow-hidden rounded-xl border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12">Select</TableHead>
              {contactColumns.map((column) => (
                <TableHead key={column.id}>{column.header}</TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {contacts.map((contact) => {
              const checked = selectedContactIds.includes(contact.id)

              return (
                <TableRow
                  key={contact.id}
                  className={onRowClick ? "cursor-pointer" : undefined}
                  onClick={onRowClick ? () => onRowClick(contact) : undefined}
                >
                  <TableCell onClick={(event) => event.stopPropagation()}>
                    <Checkbox
                      checked={checked}
                      onChange={(event) => onToggleContact(contact.id, event.target.checked)}
                    />
                  </TableCell>
                  {contactColumns.map((column) => (
                    <TableCell key={column.id}>
                      {column.id === "status" ? (
                        <Badge variant={getStatusVariant(contact.status)}>
                          {column.cell(contact)}
                        </Badge>
                      ) : column.id === "campaigns" ? (
                        contact.campaigns && contact.campaigns.length > 0 ? (
                          <button
                            type="button"
                            className="text-primary text-left underline-offset-4 hover:underline"
                            onClick={(event) => {
                              event.stopPropagation()
                              setCampaignDetailsContact(contact)
                            }}
                          >
                            {column.cell(contact)}
                          </button>
                        ) : (
                          column.cell(contact)
                        )
                      ) : (
                        column.cell(contact)
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </div>

      <Dialog
        open={Boolean(campaignDetailsContact)}
        onOpenChange={(open) => {
          if (!open) {
            setCampaignDetailsContact(null)
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Contact campaigns</DialogTitle>
            <DialogDescription>
              {campaignDetailsContact
                ? `${campaignDetailsContact.first_name} ${campaignDetailsContact.last_name}`.trim() ||
                campaignDetailsContact.email
                : "Campaign details"}
            </DialogDescription>
          </DialogHeader>
          <DialogBody className="space-y-3">
            {campaignDetailsContact?.campaigns && campaignDetailsContact.campaigns.length > 0 ? (
              campaignDetailsContact.campaigns.map((campaign) => (
                <div key={campaign.id} className="rounded-lg border p-3">
                  <p className="font-medium">{campaign.name}</p>
                  <p className="text-muted-foreground mt-1 text-xs">{campaign.id}</p>
                </div>
              ))
            ) : (
              <p className="text-muted-foreground text-sm">This contact is not attached to any campaigns.</p>
            )}
          </DialogBody>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setCampaignDetailsContact(null)}>
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
