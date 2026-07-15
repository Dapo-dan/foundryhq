import * as React from "react"

import { cn } from "@/lib/utils"
import { Input } from "@/components/ui/input"

// Not from the shadcn registry — a small local addition on top of `Input` for
// the "Show/Hide" text toggle the auth/onboarding mockups use instead of an
// eye icon.
const PasswordInput = React.forwardRef<
  HTMLInputElement,
  Omit<React.ComponentProps<"input">, "type">
>(({ className, ...props }, ref) => {
  const [visible, setVisible] = React.useState(false)

  return (
    <div className="relative">
      <Input
        ref={ref}
        type={visible ? "text" : "password"}
        className={cn("pr-14", className)}
        {...props}
      />
      <button
        type="button"
        tabIndex={-1}
        onClick={() => setVisible((v) => !v)}
        className="absolute top-1/2 right-3 -translate-y-1/2 text-xs font-medium text-brand hover:text-brand-accent"
      >
        {visible ? "Hide" : "Show"}
      </button>
    </div>
  )
})
PasswordInput.displayName = "PasswordInput"

export { PasswordInput }
