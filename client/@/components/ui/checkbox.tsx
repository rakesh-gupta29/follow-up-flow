import * as React from "react"

import { cn } from "@/lib/utils"

function Checkbox({ className, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type="checkbox"
      className={cn(
        "border-input text-primary focus-visible:ring-ring/50 size-4 rounded border bg-background shadow-xs focus-visible:ring-2 focus-visible:outline-none",
        className
      )}
      {...props}
    />
  )
}

export { Checkbox }
