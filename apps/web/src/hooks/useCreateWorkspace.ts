import { useMutation } from '@tanstack/react-query'
import { createWorkspace } from '@/services/onboarding'

export function useCreateWorkspace() {
  return useMutation({ mutationFn: createWorkspace })
}
