import { useQuery } from '@tanstack/react-query'
import { listProjects } from '@/services/projects'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useProjects() {
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  return useQuery({
    queryKey: ['projects', workspaceId],
    queryFn: listProjects,
    enabled: Boolean(workspaceId),
  })
}
