import { useForm } from "react-hook-form"
import { Navigate, useLocation, useNavigate } from "react-router-dom"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"

import { useCurrentUserQuery, useLoginMutation } from "@/lib/auth"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import type { LoginInput } from "../types/auth"

const loginInputSchema = z.object({
  email: z.email("Enter a valid email address"),
  password: z.string().min(1, "Password is required"),
})

export function LoginPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const { data: user, isLoading } = useCurrentUserQuery(true)
  const loginMutation = useLoginMutation()
  const form = useForm<LoginInput>({
    resolver: zodResolver(loginInputSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  })

  const from = (location.state as { from?: { pathname?: string } } | null)?.from?.pathname ?? "/"

  if (!isLoading && user) {
    return <Navigate to="/" replace />
  }

  const onSubmit = form.handleSubmit(async (values: LoginInput) => {
    await loginMutation.mutateAsync(values)
    navigate(from, { replace: true })
  })

  return (
    <div className="from-sidebar/40 to-background flex min-h-screen items-center justify-center bg-gradient-to-br px-4">
      <Card className="w-full max-w-md border-border/80 shadow-lg">
        <CardHeader className="space-y-2">
          <CardTitle className="text-2xl">Welcome back</CardTitle>
          <CardDescription>
            Sign in to access your logs and contacts workspace.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-5" onSubmit={onSubmit}>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="name@company.com"
                aria-invalid={Boolean(form.formState.errors.email)}
                {...form.register("email")}
              />
              {form.formState.errors.email ? (
                <p className="text-destructive text-sm">{form.formState.errors.email.message}</p>
              ) : null}
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                autoComplete="current-password"
                placeholder="Enter your password"
                aria-invalid={Boolean(form.formState.errors.password)}
                {...form.register("password")}
              />
              {form.formState.errors.password ? (
                <p className="text-destructive text-sm">{form.formState.errors.password.message}</p>
              ) : null}
            </div>

            {loginMutation.isError ? (
              <p className="text-destructive text-sm">
                {loginMutation.error instanceof Error
                  ? loginMutation.error.message
                  : "Unable to sign in"}
              </p>
            ) : null}

            <Button className="w-full" type="submit" disabled={loginMutation.isPending}>
              {loginMutation.isPending ? "Signing in..." : "Sign in"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
