import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Check, ChevronsUpDown, Plus } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { CreateWorkspaceDialog } from '@/components/layout/CreateWorkspaceDialog'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function WorkspaceSwitcher() {
  const navigate = useNavigate()
  const workspaces = useWorkspaceStore((state) => state.workspaces)
  const currentWorkspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  const setCurrentWorkspaceId = useWorkspaceStore((state) => state.setCurrentWorkspaceId)
  const currentWorkspace = workspaces.find((w) => w.id === currentWorkspaceId)
  const [createOpen, setCreateOpen] = useState(false)

  function handleSelect(id: string) {
    if (id !== currentWorkspaceId) {
      setCurrentWorkspaceId(id)
      // Workspace-scoped detail routes (e.g. sprints/:sprintId) have no
      // workspace segment in the URL, so an id from the old workspace may
      // not exist in the new one — send the user somewhere always valid.
      navigate('/dashboard')
    }
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger className="flex items-center gap-1.5 rounded-md px-2 py-1 text-sm font-medium text-muted-foreground outline-hidden hover:bg-accent hover:text-accent-foreground focus-visible:ring-3 focus-visible:ring-ring/50">
          {currentWorkspace?.name ?? 'Workspace'}
          <ChevronsUpDown size={14} className="text-muted-foreground" />
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start">
          {workspaces.map((workspace) => (
            <DropdownMenuItem key={workspace.id} onSelect={() => handleSelect(workspace.id)}>
              <span className="flex-1 truncate">{workspace.name}</span>
              {workspace.id === currentWorkspaceId && <Check size={14} />}
            </DropdownMenuItem>
          ))}
          <DropdownMenuSeparator />
          <DropdownMenuItem onSelect={() => setCreateOpen(true)}>
            <Plus size={14} />
            Create workspace
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <CreateWorkspaceDialog open={createOpen} onOpenChange={setCreateOpen} />
    </>
  )
}
