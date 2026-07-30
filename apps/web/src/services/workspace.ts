import { apiGet } from '@/lib/api-client'
import type { Workspace } from '@foundryhq/shared-types'

export function listWorkspaces() {
  return apiGet<Workspace[]>('/workspaces')
}
