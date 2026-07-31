import { useQuery } from '@tanstack/react-query'
import { listWorkspaceMembers } from '@/services/workspace'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useWorkspaceMembers() {
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  return useQuery({
    queryKey: ['workspace-members', workspaceId],
    queryFn: () => listWorkspaceMembers(workspaceId!),
    enabled: Boolean(workspaceId),
  })
}
