import { useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "../../../lib/api-client"
import type { ApiResponse } from "../../../types/api"
import type { Contact, ContactStatus } from "../../../types/contact"

export type AddContactInput = {
  email: string
  first_name: string
  last_name: string
  phone?: string
  company?: string
  status: ContactStatus
}

export async function addContact(input: AddContactInput) {
  const response = await apiClient.post<ApiResponse<Contact>>("/api/v1/add-contact", input)
  return response.data.data
}

export function useAddContactMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: addContact,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["contacts"] })
    },
  })
}
