import { CircleUserRound } from 'lucide-react'

export function TopNav() {
  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-border bg-background px-6">
      <div className="text-sm font-medium text-muted-foreground">Workspace</div>
      <CircleUserRound size={20} className="text-muted-foreground" />
    </header>
  )
}
