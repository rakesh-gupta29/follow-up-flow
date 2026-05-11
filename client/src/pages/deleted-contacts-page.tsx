import { useState } from "react"
import { useNavigate } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Select } from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useDeletedContactsQuery } from "../features/contacts/api/get-deleted-contacts"

const limits = [25, 50, 100]

export function DeletedContactsPage() {
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const [limit, setLimit] = useState(25)
  const [search, setSearch] = useState("")
  const { data, isLoading, isError, error } = useDeletedContactsQuery({
    page,
    limit,
    search,
  })

  const contacts = data?.items ?? []
  const currentPage = data?.pagination.page ?? 1
  const totalItems = data?.pagination.total ?? 0

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between gap-4">
          <div>
            <CardTitle>Deleted contacts</CardTitle>
            <p className="text-muted-foreground text-sm">
              {isLoading ? "Loading deleted contacts..." : `${totalItems} deleted contacts`}
            </p>
          </div>
          <Button type="button" variant="outline" onClick={() => navigate("/contacts")}>
            Back to contacts
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <Input
          className="w-full max-w-sm"
          placeholder="Search deleted contacts"
          value={search}
          onChange={(event) => {
            setPage(1)
            setSearch(event.target.value)
          }}
        />

        {isLoading ? (
          <div className="text-muted-foreground rounded-xl border border-dashed p-8 text-sm">
            Loading deleted contacts...
          </div>
        ) : isError ? (
          <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-4 text-sm text-destructive">
            {error instanceof Error ? error.message : "Failed to load deleted contacts."}
          </div>
        ) : (
          <div className="overflow-hidden rounded-xl border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Company</TableHead>
                  <TableHead>Phone</TableHead>
                  <TableHead>Questionnaire</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {contacts.map((contact) => (
                  <TableRow key={contact.id}>
                    <TableCell>
                      {`${contact.first_name} ${contact.last_name}`.trim() || contact.email}
                    </TableCell>
                    <TableCell>{contact.email}</TableCell>
                    <TableCell>{contact.property_name || "-"}</TableCell>
                    <TableCell>{contact.phone || "-"}</TableCell>
                    <TableCell>
                      {contact.questionnaire_url ? (
                        <a
                          href={contact.questionnaire_url}
                          target="_blank"
                          rel="noreferrer"
                          className="text-primary underline-offset-4 hover:underline"
                        >
                          Open
                        </a>
                      ) : (
                        "-"
                      )}
                    </TableCell>
                    <TableCell>{contact.status}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}

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
  )
}
