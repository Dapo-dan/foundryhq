import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createSprint } from '@/services/sprints'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useCreateSprint() {
  const queryClient = useQueryClient()
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)

  return useMutation({
    mutationFn: createSprint,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sprints', workspaceId] })
    },
  })
}
