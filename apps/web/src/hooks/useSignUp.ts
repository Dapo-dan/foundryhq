import { useMutation } from '@tanstack/react-query'
import { signUp } from '@/services/auth'
import { useAuthStore } from '@/store/slices/auth'

export function useSignUp() {
  return useMutation({
    mutationFn: signUp,
    onSuccess: (session) => useAuthStore.getState().setSession(session),
  })
}
