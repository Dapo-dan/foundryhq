import { useMutation, useQueryClient } from '@tanstack/react-query'
import { deleteTask } from '@/services/tasks'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useDeleteTask() {
  const queryClient = useQueryClient()
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)

  return useMutation({
    mutationFn: deleteTask,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks', workspaceId] })
    },
  })
}
