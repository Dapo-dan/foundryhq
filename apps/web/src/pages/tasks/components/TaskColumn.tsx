import { useDroppable } from '@dnd-kit/core'
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'
import type { Project, Task, TaskStatus } from '@foundryhq/shared-types'
import { TaskCard } from './TaskCard'

interface TaskColumnProps {
  status: TaskStatus
  label: string
  tasks: Task[]
  projectsById: Record<string, Project>
  onSelectTask: (task: Task) => void
}

export function TaskColumn({ status, label, tasks, projectsById, onSelectTask }: TaskColumnProps) {
  // The column itself is the drop target (via its status id) — a card can
  // land here even when the column is empty, which a SortableContext alone
  // (only aware of existing items) can't provide.
  const { setNodeRef } = useDroppable({ id: status })

  return (
    <div ref={setNodeRef} className="flex flex-1 flex-col gap-3 rounded-lg bg-muted/30 p-3">
      <div className="flex items-center justify-between px-1">
        <h3 className="text-sm font-semibold text-foreground">{label}</h3>
        <span className="text-xs text-muted-foreground">{tasks.length}</span>
      </div>
      <SortableContext items={tasks.map((task) => task.id)} strategy={verticalListSortingStrategy}>
        <div className="flex min-h-12 flex-col gap-2">
          {tasks.map((task) => (
            <TaskCard
              key={task.id}
              task={task}
              projectName={projectsById[task.projectId]?.name}
              onSelect={onSelectTask}
            />
          ))}
        </div>
      </SortableContext>
    </div>
  )
}
