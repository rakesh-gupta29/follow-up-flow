import { useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "../../../lib/api-client"
import type { ApiResponse } from "../../../types/api"
import type { ContactListItem } from "../../../types/contact"

export type RemoveContactsFromCampaignInput = {
  campaignId: string
  contactIds: string[]
}

export async function removeContactsFromCampaign(
  input: RemoveContactsFromCampaignInput
) {
  const results = await Promise.all(
    input.contactIds.map(async (contactId) => {
      const response = await apiClient.delete<ApiResponse<ContactListItem>>(
        `/api/v1/campaigns/${input.campaignId}/contacts/${contactId}`
      )

      return response.data.data
    })
  )

  return results
}

export function useRemoveContactsFromCampaignMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: removeContactsFromCampaign,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["contacts"] })
    },
  })
}
