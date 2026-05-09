import { useQuery } from "@tanstack/react-query"

import { apiClient } from "../../../lib/api-client"
import type { ApiResponse, PaginatedResponse } from "../../../types/api"
import {
  type ContactListItem,
  type ContactStatus,
} from "../../../types/contact"

export type ContactsQueryParams = {
  page: number
  limit: number
  search: string
  status?: ContactStatus
  campaignId?: string
}

export async function getContacts(params: ContactsQueryParams) {
  const response = await apiClient.get<ApiResponse<PaginatedResponse<ContactListItem>>>("/api/v1/contacts", {
    params: {
      page: params.page,
      limit: params.limit,
      search: params.search || undefined,
      status: params.status,
      campaignId: params.campaignId,
    },
  })
  return response.data.data
}

export function useContactsQuery(params: ContactsQueryParams) {
  return useQuery({
    queryKey: ["contacts", params],
    queryFn: () => getContacts(params),
  })
}

export type ContactsPageData = {
  items: ContactListItem[]
  page: number
  limit: number
  total: number
  total_pages: number
}
