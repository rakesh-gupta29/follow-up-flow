import { useQuery } from "@tanstack/react-query"

import { apiClient } from "../../../lib/api-client"
import type { ApiResponse } from "../../../types/api"
import type { Campaign } from "../../../types/campaign"

export async function getCampaigns() {
  const response = await apiClient.get<ApiResponse<Campaign[]>>("/api/v1/campaigns")
  return response.data.data
}

export function useCampaignsQuery() {
  return useQuery({
    queryKey: ["campaigns"],
    queryFn: getCampaigns,
  })
}
