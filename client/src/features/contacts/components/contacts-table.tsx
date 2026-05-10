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
import { campaignContactColumns, contactColumns } from "./contact-columns"
import type {
  ContactCampaign,
  ContactCampaignMembership,
  CampaignContactStatus,
  ContactListItem,
} from "../../../types/contact"

type ContactsTableProps = {
  contacts: ContactListItem[]
  selectedContactIds: string[]
  onToggleContact: (contactId: string, checked: boolean) => void
  onRowClick?: (contact: ContactListItem) => void
  variant?: "default" | "campaign"
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

function getCampaignStatusVariant(status: CampaignContactStatus) {
  if (status === "queued") {
    return "secondary"
  }

  if (status === "in_progress") {
    return "default"
  }

  if (status === "success") {
    return "outline"
  }

  return "outline"
}

function getCurrentQueuedMembership(contact: ContactListItem) {
  const queuedMemberships =
    contact.campaign_memberships?.filter(
      (membership) => membership.status === "queued"
    ) ?? []

  if (queuedMemberships.length === 0) {
    return null
  }

  return queuedMemberships.reduce((currentHighest, membership) => {
    const currentStage = currentHighest.campaign?.stage ?? -1
    const nextStage = membership.campaign?.stage ?? -1

    if (nextStage > currentStage) {
      return membership
    }

    return currentHighest
  })
}

export function ContactsTable({
  contacts,
  selectedContactIds,
  onToggleContact,
  onRowClick,
  variant = "default",
}: ContactsTableProps) {
  const [campaignDetailsContact, setCampaignDetailsContact] = useState<ContactListItem | null>(null)
  const [logsContact, setLogsContact] = useState<ContactListItem | null>(null)
  const columns = variant === "campaign" ? campaignContactColumns : contactColumns

  const renderCampaignCell = (campaigns: ContactCampaign[] | undefined, contact: ContactListItem) => {
    if (!campaigns || campaigns.length === 0) {
      return "Unassigned"
    }

    return (
      <button
        type="button"
        className="text-primary text-left underline-offset-4 hover:underline"
        onClick={(event) => {
          event.stopPropagation()
          setCampaignDetailsContact(contact)
        }}
      >
        {campaigns.length} campaign{campaigns.length === 1 ? "" : "s"}
      </button>
    )
  }

  const renderLogsButton = (contact: ContactListItem) => (
    <Button
      type="button"
      variant="outline"
      className="h-8"
      onClick={(event) => {
        event.stopPropagation()
        setLogsContact(contact)
      }}
    >
      View logs
    </Button>
  )

  const resolveMembershipName = (
    membership: ContactCampaignMembership,
    contact: ContactListItem
  ) => {
    if (membership.name) {
      return membership.name
    }

    if (membership.campaign?.name) {
      return membership.campaign.name
    }

    const matchingCampaign = contact.campaigns?.find(
      (campaign) => campaign.id === membership.campaign_id
    )

    return matchingCampaign?.name || "Unknown campaign"
  }


  return (
    <>
      <div className="overflow-hidden rounded-xl border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12">Select</TableHead>
              {columns.map((column) => (
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
                  {columns.map((column) => (
                    <TableCell key={column.id}>
                      {column.id === "status" ? (
                        <Badge variant={getStatusVariant(contact.status)}>
                          {column.cell(contact)}
                        </Badge>
                      ) : column.id === "current-campaign" ? (
                        getCurrentQueuedMembership(contact)?.campaign?.name || "-"
                      ) : column.id === "campaigns" ? (
                        renderCampaignCell(contact.campaigns, contact)
                      ) : column.id === "current-state" ? (
                        getCurrentQueuedMembership(contact) ? (
                          <Badge
                            variant={getCampaignStatusVariant(
                              getCurrentQueuedMembership(contact)!.status
                            )}
                          >
                            {getCurrentQueuedMembership(contact)!.status.replaceAll("_", " ")}
                          </Badge>
                        ) : (
                          "-"
                        )
                      ) : column.id === "logs" ? (
                        renderLogsButton(contact)
                      ) : column.id === "questionnaire_url" ? (
                        contact.questionnaire_url ? (
                          <a
                            href={contact.questionnaire_url}
                            target="_blank"
                            rel="noreferrer"
                            className="text-primary underline-offset-4 hover:underline"
                            onClick={(event) => event.stopPropagation()}
                          >
                            Open
                          </a>
                        ) : (
                          "-"
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
                  <p className="text-muted-foreground mt-1 text-xs capitalize">
                    Status: {campaign.status.replaceAll("_", " ")}
                  </p>
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

      <Dialog
        open={Boolean(logsContact)}
        onOpenChange={(open) => {
          if (!open) {
            setLogsContact(null)
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Campaign memberships</DialogTitle>
            <DialogDescription>
              {logsContact
                ? `${logsContact.first_name} ${logsContact.last_name}`.trim() || logsContact.email
                : "Campaign memberships"}
            </DialogDescription>
          </DialogHeader>
          <DialogBody className="space-y-3">
            {logsContact?.campaign_memberships &&
              logsContact.campaign_memberships.length > 0 ? (
              logsContact.campaign_memberships.map((membership, index) => {
                return <div
                  key={`${membership.campaign_id}-${membership.created_at}-${index}`}
                  className="rounded-lg border p-3"
                >
                  <div className="flex items-center justify-between gap-4">
                    <p className="font-medium">
                      {resolveMembershipName(membership, logsContact)}
                    </p>
                    <Badge variant="secondary">
                      {membership.status.replaceAll("_", " ")}
                    </Badge>
                  </div>
                  <div className="text-muted-foreground mt-2 space-y-1 text-xs">
                    <p>Added: {new Date(membership.created_at).toLocaleString()}</p>
                    <p>Updated: {new Date(membership.updated_at).toLocaleString()}</p>
                  </div>
                </div>
              })
            ) : (
              <p className="text-muted-foreground text-sm">
                No campaign memberships available for this contact.
              </p>
            )}
          </DialogBody>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setLogsContact(null)}>
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
