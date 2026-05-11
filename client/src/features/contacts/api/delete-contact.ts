import { useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "../../../lib/api-client"
import type { ApiResponse } from "../../../types/api"
import type { Contact } from "../../../types/contact"

export async function deleteContact(contactId: string) {
  const response = await apiClient.delete<ApiResponse<Contact>>(`/api/v1/contact/${contactId}`)
  return response.data.data
}

export function useDeleteContactMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: deleteContact,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["contacts"] })
      await queryClient.invalidateQueries({ queryKey: ["deleted-contacts"] })
    },
  })
}
