import { apiPost } from '@/lib/api-client'
import type { Workspace } from '@foundryhq/shared-types'

export function createWorkspace(input: { name: string }) {
  return apiPost<Workspace>('/workspaces', input)
}

// docs/api.md's invite endpoint takes one email per call
// (`POST /workspaces/{id}/members/invite`) — there's no bulk variant, so
// this fans out client-side, one request per address.
export function sendInvites(workspaceId: string, emails: string[]) {
  return Promise.all(
    emails.map((email) =>
      apiPost<void>(`/workspaces/${workspaceId}/members/invite`, { email })
    )
  ).then(() => undefined)
}
