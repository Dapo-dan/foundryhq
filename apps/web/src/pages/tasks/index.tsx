import { useMemo, useState } from 'react'
import { CheckSquare, Plus } from 'lucide-react'
import {
  DndContext,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from '@dnd-kit/core'
import type { Project, Task, TaskStatus } from '@foundryhq/shared-types'
import { Button } from '@/components/ui/button'
import { EmptyState } from '@/components/ui/empty-state'
import { useProjects } from '@/hooks/useProjects'
import { useTasks } from '@/hooks/useTasks'
import { useUpdateTask } from '@/hooks/useUpdateTask'
import { NewTaskDialog } from './components/NewTaskDialog'
import { TaskColumn } from './components/TaskColumn'

const COLUMNS: { status: TaskStatus; label: string }[] = [
  { status: 'todo', label: 'Todo' },
  { status: 'in_progress', label: 'In Progress' },
  { status: 'done', label: 'Done' },
]

export function TasksPage() {
  const [dialogOpen, setDialogOpen] = useState(false)
  const [projectFilter, setProjectFilter] = useState('')

  const { data: projects } = useProjects()
  const { data: tasks, isPending } = useTasks(projectFilter ? { projectId: projectFilter } : {})
  const updateTask = useUpdateTask()

  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 4 } }))

  const tasksById = useMemo(() => {
    const map: Record<string, Task> = {}
    for (const task of tasks ?? []) map[task.id] = task
    return map
  }, [tasks])

  const projectsById = useMemo(() => {
    const map: Record<string, Project> = {}
    for (const project of projects ?? []) map[project.id] = project
    return map
  }, [projects])

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event
    if (!over) return

    const activeTask = tasksById[active.id as string]
    if (!activeTask) return

    // Dropped directly on a column (its status id) or on another card
    // (whose status the dragged card should adopt).
    const overTask = tasksById[over.id as string]
    const newStatus = (overTask ? overTask.status : over.id) as TaskStatus

    if (newStatus !== activeTask.status) {
      updateTask.mutate({ taskId: activeTask.id, status: newStatus })
    }
  }

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
        <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
          <div className="flex gap-4">
            {COLUMNS.map((column) => (
              <TaskColumn
                key={column.status}
                status={column.status}
                label={column.label}
                tasks={(tasks ?? []).filter((task) => task.status === column.status)}
                projectsById={projectsById}
              />
            ))}
          </div>
        </DndContext>
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
