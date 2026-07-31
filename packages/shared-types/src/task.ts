export type TaskStatus = 'todo' | 'in_progress' | 'done'
export type TaskPriority = 'urgent' | 'high' | 'medium' | 'low'

export interface Task {
  id: string
  workspaceId: string
  projectId: string
  title: string
  status: TaskStatus
  assigneeId: string | null
  sprintId: string | null
  priority: TaskPriority
  storyPoints: number | null
  dueDate: string | null
  createdAt: string
  updatedAt: string
}
