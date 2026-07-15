import { useMutation } from '@tanstack/react-query'
import { sendInvites } from '@/services/onboarding'

interface SendInvitesVariables {
  workspaceId: string
  emails: string[]
}

export function useSendInvites() {
  return useMutation({
    mutationFn: ({ workspaceId, emails }: SendInvitesVariables) => sendInvites(workspaceId, emails),
  })
}
