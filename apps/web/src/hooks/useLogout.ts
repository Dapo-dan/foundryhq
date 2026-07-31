import { useMutation } from '@tanstack/react-query'
import { logout } from '@/services/auth'
import { useAuthStore } from '@/store/slices/auth'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useLogout() {
  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
      useAuthStore.getState().clearSession()
      // Otherwise a different user signing in on the same tab would
      // momentarily see the previous user's workspace name/data.
      useWorkspaceStore.getState().clear()
    },
  })
}
