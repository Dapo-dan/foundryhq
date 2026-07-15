import { useMutation } from '@tanstack/react-query'
import { forgotPassword } from '@/services/auth'

export function useForgotPassword() {
  return useMutation({ mutationFn: forgotPassword })
}
