import * as React from "react"

import { cn } from "@/lib/utils"

type DialogContextValue = {
  open: boolean
  onOpenChange: (open: boolean) => void
}

const DialogContext = React.createContext<DialogContextValue | null>(null)

type DialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  children: React.ReactNode
}

function Dialog({ open, onOpenChange, children }: DialogProps) {
  return (
    <DialogContext.Provider value={{ open, onOpenChange }}>
      {children}
    </DialogContext.Provider>
  )
}

function useDialogContext() {
  const context = React.useContext(DialogContext)

  if (!context) {
    throw new Error("Dialog components must be used inside Dialog")
  }

  return context
}

function DialogContent({
  className,
  children,
  ...props
}: React.ComponentProps<"div">) {
  const { open, onOpenChange } = useDialogContext()

  if (!open) {
    return null
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button
        aria-label="Close dialog"
        className="absolute inset-0 bg-black/50"
        onClick={() => onOpenChange(false)}
        type="button"
      />
      <div
        className={cn(
          "bg-background relative z-10 w-full max-w-lg rounded-xl border shadow-lg",
          className
        )}
        {...props}
      >
        {children}
      </div>
    </div>
  )
}

function DialogHeader({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("flex flex-col gap-2 p-6 pb-0", className)} {...props} />
}

function DialogTitle({ className, ...props }: React.ComponentProps<"h2">) {
  return <h2 className={cn("text-lg font-semibold", className)} {...props} />
}

function DialogDescription({
  className,
  ...props
}: React.ComponentProps<"p">) {
  return <p className={cn("text-muted-foreground text-sm", className)} {...props} />
}

function DialogBody({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("p-6", className)} {...props} />
}

function DialogFooter({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("flex items-center justify-end gap-2 border-t px-6 py-4", className)}
      {...props}
    />
  )
}

export {
  Dialog,
  DialogBody,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
}
