export type TaskStatus = 'todo' | 'in_progress' | 'done'

export interface Task {
  id: string
  workspaceId: string
  projectId: string
  title: string
  status: TaskStatus
  assigneeId: string | null
  createdAt: string
  updatedAt: string
}
