import { useSortable } from '@dnd-kit/sortable'
import { format } from 'date-fns'
import { Badge, type badgeVariants } from '@/components/ui/badge'
import type { Task, TaskPriority } from '@foundryhq/shared-types'
import type { VariantProps } from 'class-variance-authority'

const PRIORITY_VARIANT: Record<TaskPriority, VariantProps<typeof badgeVariants>['variant']> = {
  urgent: 'destructive',
  high: 'default',
  medium: 'secondary',
  low: 'outline',
}

interface TaskCardProps {
  task: Task
  projectName: string | undefined
  onSelect: (task: Task) => void
}

export function TaskCard({ task, projectName, onSelect }: TaskCardProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: task.id,
  })

  return (
    <div
      ref={setNodeRef}
      // @dnd-kit/utilities isn't a dependency here — building the transform
      // string by hand instead of via its CSS.Transform helper.
      style={{
        transform: transform ? `translate3d(${transform.x}px, ${transform.y}px, 0)` : undefined,
        transition,
        opacity: isDragging ? 0.4 : 1,
      }}
      {...attributes}
      {...listeners}
      // dnd-kit's PointerSensor only starts a drag past its activation
      // distance (see KanbanBoard's sensor config) — a plain click/tap never
      // crosses that threshold, so this fires normally for a real click and
      // is a no-op mid-drag.
      onClick={() => onSelect(task)}
      className="flex cursor-grab flex-col gap-2 rounded-lg border border-border bg-card p-3 text-sm active:cursor-grabbing"
    >
      <p className="font-medium text-foreground">{task.title}</p>
      <div className="flex flex-wrap items-center gap-1.5">
        {projectName && <Badge variant="secondary">{projectName}</Badge>}
        <Badge variant={PRIORITY_VARIANT[task.priority]}>{task.priority}</Badge>
        {task.storyPoints != null && <Badge variant="outline">{task.storyPoints} pts</Badge>}
      </div>
      {task.dueDate && (
        <p className="text-xs text-muted-foreground">
          Due {format(new Date(`${task.dueDate}T00:00:00`), 'MMM d')}
        </p>
      )}
    </div>
  )
}
