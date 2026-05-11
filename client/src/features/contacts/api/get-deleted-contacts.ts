import { useQuery } from "@tanstack/react-query"

import { apiClient } from "../../../lib/api-client"
import type { ApiResponse, PaginatedResponse } from "../../../types/api"
import type { Contact } from "../../../types/contact"

export type DeletedContactsQueryParams = {
  page: number
  limit: number
  search: string
}

export async function getDeletedContacts(params: DeletedContactsQueryParams) {
  const response = await apiClient.get<ApiResponse<PaginatedResponse<Contact>>>(
    "/api/v1/deleted-contacts",
    {
      params: {
        page: params.page,
        limit: params.limit,
        search: params.search || undefined,
      },
    }
  )

  return response.data.data
}

export function useDeletedContactsQuery(params: DeletedContactsQueryParams) {
  return useQuery({
    queryKey: ["deleted-contacts", params],
    queryFn: () => getDeletedContacts(params),
  })
}
