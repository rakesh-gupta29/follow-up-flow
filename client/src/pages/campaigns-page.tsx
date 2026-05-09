import { useNavigate } from "react-router-dom"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useCampaignsQuery } from "../features/campaigns/api/get-campaigns"
import { CampaignsTable } from "../features/campaigns/components/campaigns-table"

export function CampaignsPage() {
  const navigate = useNavigate()
  const { data: campaigns = [], isLoading } = useCampaignsQuery()

  return (
    <Card>
      <CardHeader>
        <CardTitle>Campaigns</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <CampaignsTable
          campaigns={campaigns}
          onRowClick={(campaign) => navigate(`/campaigns/${campaign.id}/contacts`)}
        />
        <p className="text-muted-foreground text-sm">
          {isLoading ? "Loading campaigns..." : `${campaigns.length} campaigns`}
        </p>
      </CardContent>
    </Card>
  )
}
