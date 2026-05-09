import {
  Table,
  TableBody,
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
              {campaignColumns.map((column) => column.cell(campaign))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
