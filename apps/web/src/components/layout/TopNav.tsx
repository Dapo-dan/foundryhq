import { CircleUserRound } from 'lucide-react'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function TopNav() {
  const currentWorkspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  const workspaces = useWorkspaceStore((state) => state.workspaces)
  const currentWorkspace = workspaces.find((w) => w.id === currentWorkspaceId)

  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-border bg-background px-6">
      <div className="text-sm font-medium text-muted-foreground">
        {currentWorkspace?.name ?? 'Workspace'}
      </div>
      <CircleUserRound size={20} className="text-muted-foreground" />
    </header>
  )
}
