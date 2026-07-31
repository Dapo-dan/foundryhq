import { useMemo } from 'react'
import { format } from 'date-fns'
import { ArrowLeft } from 'lucide-react'
import { Link, useParams } from 'react-router-dom'
import type { Project } from '@foundryhq/shared-types'
import { KanbanBoard } from '@/pages/tasks/components/KanbanBoard'
import { useProjects } from '@/hooks/useProjects'
import { useSprintDetail } from '@/hooks/useSprintDetail'
import { useTasks } from '@/hooks/useTasks'
import { useVelocity } from '@/hooks/useVelocity'

function formatDate(dateOnly: string): string {
  return format(new Date(`${dateOnly}T00:00:00`), 'MMM d, yyyy')
}

export function SprintDetailPage() {
  const { sprintId } = useParams<{ sprintId: string }>()
  const { data, isPending } = useSprintDetail(sprintId)
  const { data: velocity } = useVelocity(sprintId)
  const { data: projects } = useProjects()
  // Sourced from the shared 'tasks' query (filtered by sprintId) rather than
  // the sprint bundle's own tasks — useUpdateTask's optimistic cache patch
  // (ADR-0005) only targets ['tasks', workspaceId] queries, so the Kanban
  // board here needs to read from that same cache to see live drag updates.
  const { data: tasks } = useTasks(sprintId ? { sprintId } : {})

  const projectsById = useMemo(() => {
    const map: Record<string, Project> = {}
    for (const project of projects ?? []) map[project.id] = project
    return map
  }, [projects])

  if (isPending) {
    return <p className="text-sm text-muted-foreground">Loading sprint…</p>
  }

  if (!data) {
    return <p className="text-sm text-muted-foreground">Sprint not found.</p>
  }

  return (
    <div className="flex flex-col gap-6">
      <div>
        <Link
          to="/sprints"
          className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft size={16} />
          Back to sprints
        </Link>
        <div className="mt-2 flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold">{data.sprint.name}</h1>
            <p className="text-sm text-muted-foreground">
              {formatDate(data.sprint.startDate)} – {formatDate(data.sprint.endDate)}
            </p>
          </div>
          <div className="text-right">
            <p className="text-2xl font-semibold">{(velocity?.velocity ?? 0).toFixed(1)}</p>
            <p className="text-xs text-muted-foreground">velocity (story points done)</p>
          </div>
        </div>
      </div>

      <KanbanBoard tasks={tasks ?? []} projectsById={projectsById} />
    </div>
  )
}
