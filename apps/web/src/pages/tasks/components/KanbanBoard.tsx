import { useMemo } from 'react'
import {
  DndContext,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from '@dnd-kit/core'
import type { Project, Task, TaskStatus } from '@foundryhq/shared-types'
import { useUpdateTask } from '@/hooks/useUpdateTask'
import { TaskColumn } from './TaskColumn'

const COLUMNS: { status: TaskStatus; label: string }[] = [
  { status: 'todo', label: 'Todo' },
  { status: 'in_progress', label: 'In Progress' },
  { status: 'done', label: 'Done' },
]

interface KanbanBoardProps {
  tasks: Task[]
  projectsById: Record<string, Project>
}

// Extracted from TasksPage so the Sprints detail page can show the same
// drag-and-drop board for a sprint's tasks instead of duplicating it.
export function KanbanBoard({ tasks, projectsById }: KanbanBoardProps) {
  const updateTask = useUpdateTask()
  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 4 } }))

  const tasksById = useMemo(() => {
    const map: Record<string, Task> = {}
    for (const task of tasks) map[task.id] = task
    return map
  }, [tasks])

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

  return (
    <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
      <div className="flex gap-4">
        {COLUMNS.map((column) => (
          <TaskColumn
            key={column.status}
            status={column.status}
            label={column.label}
            tasks={tasks.filter((task) => task.status === column.status)}
            projectsById={projectsById}
          />
        ))}
      </div>
    </DndContext>
  )
}
