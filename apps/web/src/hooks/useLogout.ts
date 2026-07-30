import { useMutation } from '@tanstack/react-query'
import { logout } from '@/services/auth'
import { useAuthStore } from '@/store/slices/auth'

export function useLogout() {
  return useMutation({
    mutationFn: logout,
    onSuccess: () => useAuthStore.getState().clearSession(),
  })
}
