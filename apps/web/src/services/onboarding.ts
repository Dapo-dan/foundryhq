import { apiPost } from '@/lib/api-client'
import type { CreateWorkspaceInput, Workspace } from '@/types/onboarding'

// NOTE: docs/api.md only documents GET/PATCH /workspaces/{id} — creation
// isn't listed. Guessed as a standard REST POST to the collection; confirm
// against the real handler once it exists.
export function createWorkspace(input: CreateWorkspaceInput) {
  return apiPost<Workspace>('/workspaces', input)
}

// NOTE: docs/api.md's invite endpoint takes one email per call
// (`POST /workspaces/{id}/members/invite`) — there's no documented bulk
// variant, so this fans out client-side, one request per address.
export function sendInvites(workspaceId: string, emails: string[]) {
  return Promise.all(
    emails.map((email) =>
      apiPost<void>(`/workspaces/${workspaceId}/members/invite`, { email })
    )
  ).then(() => undefined)
}
