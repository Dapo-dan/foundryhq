import { useQuery } from '@tanstack/react-query'
import { getSprintVelocity } from '@/services/sprints'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useVelocity(sprintId: string | undefined) {
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  return useQuery({
    queryKey: ['sprints', workspaceId, sprintId, 'velocity'],
    queryFn: () => getSprintVelocity(sprintId!),
    enabled: Boolean(workspaceId) && Boolean(sprintId),
  })
}
