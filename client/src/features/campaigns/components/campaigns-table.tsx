import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { campaignColumns } from "./campaign-columns"
import type { Campaign } from "../../../types/campaign"

type CampaignsTableProps = {
  campaigns: Campaign[]
  onRowClick: (campaign: Campaign) => void
}

function getStatusVariant(status: Campaign["status"]) {
  if (status === "active") {
    return "default"
  }

  if (status === "paused") {
    return "secondary"
  }

  return "outline"
}

export function CampaignsTable({ campaigns, onRowClick }: CampaignsTableProps) {
  return (
    <div className="overflow-hidden rounded-xl border">
      <Table>
        <TableHeader>
          <TableRow>
            {campaignColumns.map((column) => (
              <TableHead key={column.id}>{column.header}</TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {campaigns.map((campaign) => (
            <TableRow
              key={campaign.id}
              className="cursor-pointer"
              onClick={() => onRowClick(campaign)}
            >
              {campaignColumns.map((column) => (
                <TableCell key={column.id}>
                  {column.id === "status" ? (
                    <Badge variant={getStatusVariant(campaign.status)}>
                      {column.cell(campaign)}
                    </Badge>
                  ) : (
                    column.cell(campaign)
                  )}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
