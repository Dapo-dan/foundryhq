import { useMutation } from '@tanstack/react-query'
import { acceptInvite } from '@/services/auth'
import { useAuthStore } from '@/store/slices/auth'

export function useAcceptInvite() {
  return useMutation({
    mutationFn: acceptInvite,
    onSuccess: (session) => useAuthStore.getState().setSession(session),
  })
}
