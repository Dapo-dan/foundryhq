import { apiDelete, apiGet, apiPatch, apiPost } from '@/lib/api-client'
import type { Task, TaskPriority, TaskStatus } from '@foundryhq/shared-types'
import { useWorkspaceStore } from '@/store/slices/workspace'

// Reads the current workspace from the store rather than taking it as a
// parameter — same convention as services/projects.ts.
function requireWorkspaceId(): string {
  const workspaceId = useWorkspaceStore.getState().currentWorkspaceId
  if (!workspaceId) {
    throw new Error('No current workspace — cannot call a workspace-scoped endpoint yet.')
  }
  return workspaceId
}

export interface TaskFilters {
  projectId?: string
  status?: TaskStatus
  assigneeId?: string
  sprintId?: string
}

export function listTasks(filters: TaskFilters = {}) {
  const params = new URLSearchParams()
  if (filters.projectId) params.set('projectId', filters.projectId)
  if (filters.status) params.set('status', filters.status)
  if (filters.assigneeId) params.set('assigneeId', filters.assigneeId)
  if (filters.sprintId) params.set('sprintId', filters.sprintId)
  const query = params.toString()
  return apiGet<Task[]>(`/workspaces/${requireWorkspaceId()}/tasks${query ? `?${query}` : ''}`)
}

export interface CreateTaskInput {
  projectId: string
  title: string
  assigneeId?: string
  sprintId?: string
  priority?: TaskPriority
  storyPoints?: number
  dueDate?: string
}

export function createTask(input: CreateTaskInput) {
  return apiPost<Task>(`/workspaces/${requireWorkspaceId()}/tasks`, input)
}

export interface UpdateTaskInput {
  projectId?: string
  title?: string
  status?: TaskStatus
  assigneeId?: string
  clearAssignee?: boolean
  sprintId?: string
  priority?: TaskPriority
  storyPoints?: number
  clearStoryPoints?: boolean
  dueDate?: string
}

export function updateTask(taskId: string, input: UpdateTaskInput) {
  return apiPatch<Task>(`/workspaces/${requireWorkspaceId()}/tasks/${taskId}`, input)
}

export function deleteTask(taskId: string) {
  return apiDelete<void>(`/workspaces/${requireWorkspaceId()}/tasks/${taskId}`)
}
