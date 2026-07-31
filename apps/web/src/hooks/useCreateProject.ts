import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createProject } from '@/services/projects'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useCreateProject() {
  const queryClient = useQueryClient()
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)

  return useMutation({
    mutationFn: createProject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects', workspaceId] })
    },
  })
}
