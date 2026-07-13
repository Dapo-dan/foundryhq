import { FolderKanban, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function ProjectsPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Projects</h1>
          <p className="text-sm text-muted-foreground">
            Track the projects your team is working on.
          </p>
        </div>
        <Button>
          <Plus size={20} />
          New project
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <FolderKanban size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">No projects yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Create your first project to start organizing work for your team.
        </p>
        <Button variant="outline">Create project</Button>
      </div>
    </div>
  )
}
