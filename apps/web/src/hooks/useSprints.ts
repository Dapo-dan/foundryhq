import { useQuery } from '@tanstack/react-query'
import { listSprints } from '@/services/sprints'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useSprints() {
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  return useQuery({
    queryKey: ['sprints', workspaceId],
    queryFn: listSprints,
    enabled: Boolean(workspaceId),
  })
}
