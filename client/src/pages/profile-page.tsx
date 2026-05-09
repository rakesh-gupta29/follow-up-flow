import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useCurrentUserQuery } from "@/lib/auth"

export function ProfilePage() {
  const { data: user } = useCurrentUserQuery(true)

  return (
    <Card className="max-w-2xl">
      <CardHeader>
        <CardTitle>Profile</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-sm">
          Signed in as <span className="font-medium">{user?.email ?? "Unknown user"}</span>
        </p>
      </CardContent>
    </Card>
  )
}
