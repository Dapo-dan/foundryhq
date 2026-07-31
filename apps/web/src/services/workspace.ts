import { apiGet } from '@/lib/api-client'
import type { Workspace, WorkspaceMember } from '@foundryhq/shared-types'

export function listWorkspaces() {
  return apiGet<Workspace[]>('/workspaces')
}

export function listWorkspaceMembers(workspaceId: string) {
  return apiGet<WorkspaceMember[]>(`/workspaces/${workspaceId}/members`)
}
