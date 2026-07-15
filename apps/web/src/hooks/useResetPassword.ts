import { useMutation } from '@tanstack/react-query'
import { resetPassword } from '@/services/auth'

export function useResetPassword() {
  return useMutation({ mutationFn: resetPassword })
}
