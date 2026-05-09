import type { ReactNode } from "react"
import { Navigate, useLocation } from "react-router-dom"

import { useCurrentUserQuery } from "@/lib/auth"
import { clearAccessToken, getAccessToken } from "./lib/access-token"

type AuthGuardProps = {
  children: ReactNode
}

export function AuthGuard({ children }: AuthGuardProps) {
  const location = useLocation()
  const token = getAccessToken()
  const { data, isLoading, isError } = useCurrentUserQuery(true)

  if (!token) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-muted-foreground text-sm">Authenticating...</p>
      </div>
    )
  }

  if (isError || !data) {
    clearAccessToken()
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return <>{children}</>
}
