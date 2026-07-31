import { apiGet, apiPost } from '@/lib/api-client'
import type { Sprint, Task } from '@foundryhq/shared-types'
import { useWorkspaceStore } from '@/store/slices/workspace'

// Reads the current workspace from the store rather than taking it as a
// parameter — same convention as services/projects.ts and services/tasks.ts.
function requireWorkspaceId(): string {
  const workspaceId = useWorkspaceStore.getState().currentWorkspaceId
  if (!workspaceId) {
    throw new Error('No current workspace — cannot call a workspace-scoped endpoint yet.')
  }
  return workspaceId
}

export function listSprints() {
  return apiGet<Sprint[]>(`/workspaces/${requireWorkspaceId()}/sprints`)
}

export function createSprint(input: { name: string; startDate: string; endDate: string }) {
  return apiPost<Sprint>(`/workspaces/${requireWorkspaceId()}/sprints`, input)
}

export interface SprintWithTasks {
  sprint: Sprint
  tasks: Task[]
}

export function getSprint(sprintId: string) {
  return apiGet<SprintWithTasks>(`/workspaces/${requireWorkspaceId()}/sprints/${sprintId}`)
}

export interface SprintVelocity {
  velocity: number
}

export function getSprintVelocity(sprintId: string) {
  return apiGet<SprintVelocity>(`/workspaces/${requireWorkspaceId()}/sprints/${sprintId}/velocity`)
}
