import { useSortable } from '@dnd-kit/sortable'
import { Badge } from '@/components/ui/badge'
import type { Task } from '@foundryhq/shared-types'

interface TaskCardProps {
  task: Task
  projectName: string | undefined
}

export function TaskCard({ task, projectName }: TaskCardProps) {
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
      className="flex cursor-grab flex-col gap-2 rounded-lg border border-border bg-card p-3 text-sm active:cursor-grabbing"
    >
      <p className="font-medium text-foreground">{task.title}</p>
      {projectName && <Badge variant="secondary">{projectName}</Badge>}
    </div>
  )
}
