import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createTask } from '@/services/tasks'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useCreateTask() {
  const queryClient = useQueryClient()
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)

  return useMutation({
    mutationFn: createTask,
    onSuccess: () => {
      // Partial key — matches every cached filter variant of the tasks
      // list, not just an unfiltered one.
      queryClient.invalidateQueries({ queryKey: ['tasks', workspaceId] })
    },
  })
}
