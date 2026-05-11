import { useMemo, useState } from "react"
import { useNavigate } from "react-router-dom"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"

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
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { useDeleteContactMutation } from "../features/contacts/api/delete-contact"
import { useContactsQuery } from "../features/contacts/api/get-contacts"
import { useUpdateContactMutation } from "../features/contacts/api/update-contact"
import { AddToCampaignDialog } from "../features/contacts/components/add-to-campaign-dialog"
import { ContactsTable } from "../features/contacts/components/contacts-table"
import type { ContactListItem, ContactStatus } from "../types/contact"

const limits = [25, 50, 100]
const editContactSchema = z.object({
  email: z.email("Enter a valid email address"),
  first_name: z.string().min(1, "First name is required"),
  last_name: z.string().min(1, "Last name is required"),
  property_name: z.string().optional(),
  phone: z.string().optional(),
  questionnaire_url: z.string().optional(),
  thread_id: z.string().optional(),
  status: z.enum(["active", "unsubscribed", "bounced"]),
})

type EditContactFormValues = z.infer<typeof editContactSchema>

export function ContactsPage() {
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const [limit, setLimit] = useState(25)
  const [search, setSearch] = useState("")
  const [status, setStatus] = useState<ContactStatus | undefined>(undefined)
  const [selectedContactIds, setSelectedContactIds] = useState<string[]>([])
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingContact, setEditingContact] = useState<ContactListItem | null>(null)
  const [deletingContact, setDeletingContact] = useState<ContactListItem | null>(null)
  const updateContactMutation = useUpdateContactMutation()
  const deleteContactMutation = useDeleteContactMutation()
  const editForm = useForm<EditContactFormValues>({
    resolver: zodResolver(editContactSchema),
    defaultValues: {
      email: "",
      first_name: "",
      last_name: "",
      property_name: "",
      phone: "",
      questionnaire_url: "",
      thread_id: "",
      status: "active",
    },
  })
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

  const openEditDialog = (contact: ContactListItem) => {
    setEditingContact(contact)
    editForm.reset({
      email: contact.email,
      first_name: contact.first_name,
      last_name: contact.last_name,
      property_name: contact.property_name || "",
      phone: contact.phone || "",
      questionnaire_url: contact.questionnaire_url || "",
      thread_id: contact.thread_id || "",
      status: contact.status,
    })
  }

  const handleEditSubmit = editForm.handleSubmit(async (values) => {
    if (!editingContact) {
      return
    }

    await updateContactMutation.mutateAsync({
      id: editingContact.id,
      updates: {
        ...values,
        property_name: values.property_name || undefined,
        phone: values.phone || undefined,
        questionnaire_url: values.questionnaire_url || undefined,
        thread_id: values.thread_id || undefined,
      },
    })

    setEditingContact(null)
  })

  const handleDeleteContact = async () => {
    if (!deletingContact) {
      return
    }

    await deleteContactMutation.mutateAsync(deletingContact.id)
    setDeletingContact(null)
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
              <Button type="button" onClick={() => navigate("/contacts/add")}>
                Add contacts
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate("/contacts/deleted")}>
                Deleted contacts
              </Button>
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
            variant="campaign"
            renderActions={(contact) => (
              <div className="flex items-center gap-2">
                <Button type="button" variant="outline" className="h-8" onClick={() => openEditDialog(contact)}>
                  Edit
                </Button>
                <Button type="button" variant="outline" className="h-8" onClick={() => setDeletingContact(contact)}>
                  Delete
                </Button>
              </div>
            )}
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

      <Dialog open={Boolean(editingContact)} onOpenChange={(open) => !open && setEditingContact(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit contact</DialogTitle>
            <DialogDescription>Update the selected contact details.</DialogDescription>
          </DialogHeader>
          <DialogBody>
            <form className="space-y-4" onSubmit={handleEditSubmit}>
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="edit-first-name">First name</Label>
                  <Input id="edit-first-name" {...editForm.register("first_name")} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="edit-last-name">Last name</Label>
                  <Input id="edit-last-name" {...editForm.register("last_name")} />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-email">Email</Label>
                <Input id="edit-email" type="email" {...editForm.register("email")} />
              </div>
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="edit-property-name">Company</Label>
                  <Input id="edit-property-name" {...editForm.register("property_name")} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="edit-phone">Phone</Label>
                  <Input id="edit-phone" {...editForm.register("phone")} />
                </div>
              </div>
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="edit-questionnaire-url">Questionnaire URL</Label>
                  <Input id="edit-questionnaire-url" {...editForm.register("questionnaire_url")} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="edit-thread-id">Thread ID</Label>
                  <Input id="edit-thread-id" {...editForm.register("thread_id")} />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-status">Status</Label>
                <Select id="edit-status" {...editForm.register("status")}>
                  <option value="active">Active</option>
                  <option value="unsubscribed">Unsubscribed</option>
                  <option value="bounced">Bounced</option>
                </Select>
              </div>
              {updateContactMutation.isError ? (
                <p className="text-destructive text-sm">
                  {updateContactMutation.error instanceof Error
                    ? updateContactMutation.error.message
                    : "Failed to update contact"}
                </p>
              ) : null}
              <DialogFooter className="px-0 pb-0">
                <Button type="button" variant="outline" onClick={() => setEditingContact(null)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={updateContactMutation.isPending}>
                  {updateContactMutation.isPending ? "Saving..." : "Save changes"}
                </Button>
              </DialogFooter>
            </form>
          </DialogBody>
        </DialogContent>
      </Dialog>

      <Dialog open={Boolean(deletingContact)} onOpenChange={(open) => !open && setDeletingContact(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete contact</DialogTitle>
            <DialogDescription>
              This will move the contact into deleted contacts.
            </DialogDescription>
          </DialogHeader>
          <DialogBody className="text-sm">
            {deletingContact
              ? `${deletingContact.first_name} ${deletingContact.last_name}`.trim() || deletingContact.email
              : ""}
          </DialogBody>
          {deleteContactMutation.isError ? (
            <DialogBody className="pt-0 text-sm text-destructive">
              {deleteContactMutation.error instanceof Error
                ? deleteContactMutation.error.message
                : "Failed to delete contact"}
            </DialogBody>
          ) : null}
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setDeletingContact(null)}>
              Cancel
            </Button>
            <Button type="button" onClick={handleDeleteContact} disabled={deleteContactMutation.isPending}>
              {deleteContactMutation.isPending ? "Deleting..." : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
