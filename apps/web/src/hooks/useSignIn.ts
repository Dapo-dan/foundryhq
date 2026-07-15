import { useMutation } from '@tanstack/react-query'
import { signIn } from '@/services/auth'
import { useAuthStore } from '@/store/slices/auth'

export function useSignIn() {
  return useMutation({
    mutationFn: signIn,
    onSuccess: (session) => useAuthStore.getState().setSession(session),
  })
}
