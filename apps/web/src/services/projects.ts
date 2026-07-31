import { apiGet, apiPost } from '@/lib/api-client'
import type { Project } from '@foundryhq/shared-types'
import { useWorkspaceStore } from '@/store/slices/workspace'

// Reads the current workspace from the store rather than taking it as a
// parameter — every call site is already inside a signed-in, workspace-
// scoped part of the app (see AppLayout's onboarding redirect), so there's
// always exactly one to read.
function requireWorkspaceId(): string {
  const workspaceId = useWorkspaceStore.getState().currentWorkspaceId
  if (!workspaceId) {
    throw new Error('No current workspace — cannot call a workspace-scoped endpoint yet.')
  }
  return workspaceId
}

export function listProjects() {
  return apiGet<Project[]>(`/workspaces/${requireWorkspaceId()}/projects`)
}

export function createProject(input: { name: string; description?: string }) {
  return apiPost<Project>(`/workspaces/${requireWorkspaceId()}/projects`, input)
}
