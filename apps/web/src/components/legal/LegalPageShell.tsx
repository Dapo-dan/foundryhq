import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'

interface LegalPageShellProps {
  title: string
  children: ReactNode
}

export function LegalPageShell({ title, children }: LegalPageShellProps) {
  return (
    <div className="mx-auto flex min-h-svh max-w-2xl flex-col gap-8 px-6 py-16">
      <Link
        to="/"
        className="rounded-sm text-sm text-muted-foreground hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
      >
        ← Back to home
      </Link>
      <div className="flex flex-col gap-4">
        <h1 className="text-3xl font-bold">{title}</h1>
        <div className="flex flex-col gap-3 text-sm leading-relaxed text-muted-foreground">
          {children}
        </div>
      </div>
    </div>
  )
}
