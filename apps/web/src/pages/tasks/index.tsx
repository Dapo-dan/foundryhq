import { useMemo, useState } from 'react'
import { CheckSquare, Plus } from 'lucide-react'
import type { Project } from '@foundryhq/shared-types'
import { Button } from '@/components/ui/button'
import { EmptyState } from '@/components/ui/empty-state'
import { useProjects } from '@/hooks/useProjects'
import { useTasks } from '@/hooks/useTasks'
import { KanbanBoard } from './components/KanbanBoard'
import { NewTaskDialog } from './components/NewTaskDialog'

export function TasksPage() {
  const [dialogOpen, setDialogOpen] = useState(false)
  const [projectFilter, setProjectFilter] = useState('')

  const { data: projects } = useProjects()
  const { data: tasks, isPending } = useTasks(projectFilter ? { projectId: projectFilter } : {})

  const projectsById = useMemo(() => {
    const map: Record<string, Project> = {}
    for (const project of projects ?? []) map[project.id] = project
    return map
  }, [projects])

  const hasTasks = (tasks?.length ?? 0) > 0
  const hasProjects = (projects?.length ?? 0) > 0

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Tasks</h1>
          <p className="text-sm text-muted-foreground">
            See and manage the tasks assigned to your workspace.
          </p>
        </div>
        <Button onClick={() => setDialogOpen(true)}>
          <Plus size={20} />
          New task
        </Button>
      </div>

      {hasProjects && (
        <select
          value={projectFilter}
          onChange={(e) => setProjectFilter(e.target.value)}
          className="h-8 w-fit min-w-40 rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 dark:bg-input/30"
        >
          <option value="">All projects</option>
          {projects!.map((project) => (
            <option key={project.id} value={project.id}>
              {project.name}
            </option>
          ))}
        </select>
      )}

      {isPending ? (
        <p className="text-sm text-muted-foreground">Loading tasks…</p>
      ) : hasTasks ? (
        <KanbanBoard tasks={tasks ?? []} projectsById={projectsById} />
      ) : (
        <EmptyState
          icon={CheckSquare}
          title="No tasks yet"
          description="Create a task to start tracking work for your team."
          action={
            <Button variant="outline" onClick={() => setDialogOpen(true)}>
              Create task
            </Button>
          }
        />
      )}

      <NewTaskDialog open={dialogOpen} onOpenChange={setDialogOpen} projects={projects ?? []} />
    </div>
  )
}
