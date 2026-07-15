import type { ReactNode } from 'react'

interface AuthCardProps {
  heading: string
  description: string
  children: ReactNode
}

// Shared shell for the 4 auth screens (Sign In, Sign Up, Forgot/Reset
// Password) — matches the mockups' plain centered heading + subtext + form
// stack, with no visible card border (the design has no boxed container here,
// unlike the shadcn `Card` primitive's default ring/background).
export function AuthCard({ heading, description, children }: AuthCardProps) {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1 text-center">
        <h1 className="text-2xl font-bold text-foreground">{heading}</h1>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
      {children}
    </div>
  )
}
