import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { useNavigate } from "react-router-dom"
import { z } from "zod"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { useAddContactMutation } from "../features/contacts/api/add-contact"
import type { ContactStatus } from "../types/contact"

const addContactSchema = z.object({
  email: z.email("Enter a valid email address"),
  first_name: z.string().min(1, "First name is required"),
  last_name: z.string().min(1, "Last name is required"),
  phone: z.string().optional(),
  company: z.string().optional(),
  status: z.enum(["active", "unsubscribed", "bounced"]),
})

type AddContactFormValues = z.infer<typeof addContactSchema>

export function AddContactsPage() {
  const navigate = useNavigate()
  const addContactMutation = useAddContactMutation()
  const form = useForm<AddContactFormValues>({
    resolver: zodResolver(addContactSchema),
    defaultValues: {
      email: "",
      first_name: "",
      last_name: "",
      phone: "",
      company: "",
      status: "active",
    },
  })

  const onSubmit = form.handleSubmit(async (values) => {
    await addContactMutation.mutateAsync({
      ...values,
      phone: values.phone || undefined,
      company: values.company || undefined,
      status: values.status as ContactStatus,
    })
    navigate("/contacts", { replace: true })
  })

  return (
    <Card className="max-w-2xl">
      <CardHeader>
        <CardTitle>Add contacts</CardTitle>
      </CardHeader>
      <CardContent>
        <form className="space-y-5" onSubmit={onSubmit}>
          <div className="grid gap-5 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="first_name">First name</Label>
              <Input id="first_name" {...form.register("first_name")} />
              {form.formState.errors.first_name ? (
                <p className="text-destructive text-sm">{form.formState.errors.first_name.message}</p>
              ) : null}
            </div>
            <div className="space-y-2">
              <Label htmlFor="last_name">Last name</Label>
              <Input id="last_name" {...form.register("last_name")} />
              {form.formState.errors.last_name ? (
                <p className="text-destructive text-sm">{form.formState.errors.last_name.message}</p>
              ) : null}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" {...form.register("email")} />
            {form.formState.errors.email ? (
              <p className="text-destructive text-sm">{form.formState.errors.email.message}</p>
            ) : null}
          </div>

          <div className="grid gap-5 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="phone">Phone</Label>
              <Input id="phone" {...form.register("phone")} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="company">Company</Label>
              <Input id="company" {...form.register("company")} />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <Select id="status" {...form.register("status")}>
              <option value="active">Active</option>
              <option value="unsubscribed">Unsubscribed</option>
              <option value="bounced">Bounced</option>
            </Select>
          </div>

          {addContactMutation.isError ? (
            <p className="text-destructive text-sm">
              {addContactMutation.error instanceof Error
                ? addContactMutation.error.message
                : "Failed to add contact"}
            </p>
          ) : null}

          <div className="flex items-center gap-3">
            <Button type="submit" disabled={addContactMutation.isPending}>
              {addContactMutation.isPending ? "Adding..." : "Add contact"}
            </Button>
            <Button type="button" variant="outline" onClick={() => navigate("/contacts")}>
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
