import { Badge } from "@/components/ui/badge"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
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
  return (
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
  )
}
