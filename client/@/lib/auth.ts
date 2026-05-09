import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "../../src/lib/api-client"
import type { LoginInput, LoginResponse } from "../../src/types/auth"
import type { ApiResponse } from "../../src/types/api"
import type { CurrentUser } from "../../src/types/profile"
import {
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from "../../src/lib/access-token"

export async function login(input: LoginInput) {
  const response = await apiClient.post<ApiResponse<LoginResponse>>("/api/v1/auth/login", input)
  const parsed = response.data.data

  setAccessToken(parsed.access_token)

  return parsed
}

export async function fetchCurrentUser() {
  const token = getAccessToken()

  if (!token) {
    throw new Error("Missing access token")
  }

  const response = await apiClient.get<ApiResponse<CurrentUser>>("/api/v1/auth/me")

  return response.data.data
}

export function useCurrentUserQuery(enabled = true) {
  return useQuery({
    queryKey: ["auth", "me"],
    queryFn: fetchCurrentUser,
    enabled: enabled && Boolean(getAccessToken()),
    retry: false,
  })
}

export function useLoginMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: login,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["auth", "me"] })
    },
  })
}

export function logout() {
  clearAccessToken()
}
