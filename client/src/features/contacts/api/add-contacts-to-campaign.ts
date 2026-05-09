import { useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "../../../lib/api-client"
import type { ApiResponse } from "../../../types/api"

export type AddContactsToCampaignInput = {
  campaignId: string
  contactIds: string[]
}

export type AddContactsToCampaignResponse = {
  success?: boolean
}

export async function addContactsToCampaign(input: AddContactsToCampaignInput) {
  const response = await apiClient.post(
    `/api/v1/campaigns/${input.campaignId}/contacts`,
    { contact_ids: input.contactIds }
  )

  return (response.data as ApiResponse<AddContactsToCampaignResponse>).data
}

export function useAddContactsToCampaignMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: addContactsToCampaign,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["contacts"] })
    },
  })
}
