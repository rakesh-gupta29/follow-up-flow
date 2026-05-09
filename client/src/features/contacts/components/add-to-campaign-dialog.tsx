import { useEffect, useState } from "react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogBody,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Select } from "@/components/ui/select"
import { useCampaignsQuery } from "../../campaigns/api/get-campaigns"
import { useAddContactsToCampaignMutation } from "../api/add-contacts-to-campaign"

type AddToCampaignDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  selectedContactIds: string[]
}

export function AddToCampaignDialog({
  open,
  onOpenChange,
  selectedContactIds,
}: AddToCampaignDialogProps) {
  const { data: campaigns = [], isLoading } = useCampaignsQuery()
  const addToCampaignMutation = useAddContactsToCampaignMutation()
  const [campaignId, setCampaignId] = useState("")

  useEffect(() => {
    if (!open) {
      setCampaignId("")
    }
  }, [open])

  const handleSubmit = async () => {
    await addToCampaignMutation.mutateAsync({
      campaignId,
      contactIds: selectedContactIds,
    })
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add to campaign</DialogTitle>
          <DialogDescription>
            Choose a campaign for the {selectedContactIds.length} selected contacts.
          </DialogDescription>
        </DialogHeader>
        <DialogBody>
          <Select
            value={campaignId}
            onChange={(event) => setCampaignId(event.target.value)}
            disabled={isLoading || campaigns.length === 0}
            className="w-full"
          >
            <option value="">Select a campaign</option>
            {campaigns.map((campaign) => (
              <option key={campaign.id} value={campaign.id}>
                {campaign.name}
              </option>
            ))}
          </Select>
        </DialogBody>
        <DialogFooter>
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            type="button"
            onClick={handleSubmit}
            disabled={!campaignId || addToCampaignMutation.isPending}
          >
            {addToCampaignMutation.isPending ? "Adding..." : "Add to campaign"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
