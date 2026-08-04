import { CircleUserRound } from 'lucide-react'
import { WorkspaceSwitcher } from '@/components/layout/WorkspaceSwitcher'

export function TopNav() {
  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-border bg-background px-6">
      <WorkspaceSwitcher />
      <CircleUserRound size={20} className="text-muted-foreground" />
    </header>
  )
}
