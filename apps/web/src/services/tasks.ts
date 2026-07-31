import { apiDelete, apiGet, apiPatch, apiPost } from '@/lib/api-client'
import type { Task, TaskStatus } from '@foundryhq/shared-types'
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
}

export function listTasks(filters: TaskFilters = {}) {
  const params = new URLSearchParams()
  if (filters.projectId) params.set('projectId', filters.projectId)
  if (filters.status) params.set('status', filters.status)
  if (filters.assigneeId) params.set('assigneeId', filters.assigneeId)
  const query = params.toString()
  return apiGet<Task[]>(`/workspaces/${requireWorkspaceId()}/tasks${query ? `?${query}` : ''}`)
}

export function createTask(input: { projectId: string; title: string; assigneeId?: string }) {
  return apiPost<Task>(`/workspaces/${requireWorkspaceId()}/tasks`, input)
}

export interface UpdateTaskInput {
  projectId?: string
  title?: string
  status?: TaskStatus
  assigneeId?: string
  clearAssignee?: boolean
}

export function updateTask(taskId: string, input: UpdateTaskInput) {
  return apiPatch<Task>(`/workspaces/${requireWorkspaceId()}/tasks/${taskId}`, input)
}

export function deleteTask(taskId: string) {
  return apiDelete<void>(`/workspaces/${requireWorkspaceId()}/tasks/${taskId}`)
}
