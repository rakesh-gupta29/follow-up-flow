import { useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "../../../lib/api-client"
import type { ApiResponse } from "../../../types/api"
import type { Contact } from "../../../types/contact"

export type UpdateContactInput = {
  id: string
  updates: Partial<
    Pick<
      Contact,
      | "email"
      | "first_name"
      | "last_name"
      | "property_name"
      | "phone"
      | "questionnaire_url"
      | "thread_id"
      | "status"
    >
  >
}

export async function updateContact(input: UpdateContactInput) {
  const response = await apiClient.patch<ApiResponse<Contact>>(
    `/api/v1/contact/${input.id}`,
    input.updates
  )

  return response.data.data
}

export function useUpdateContactMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: updateContact,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["contacts"] })
      await queryClient.invalidateQueries({ queryKey: ["deleted-contacts"] })
    },
  })
}
