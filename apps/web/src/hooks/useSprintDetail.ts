import { useQuery } from '@tanstack/react-query'
import { getSprint } from '@/services/sprints'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useSprintDetail(sprintId: string | undefined) {
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  return useQuery({
    queryKey: ['sprints', workspaceId, sprintId],
    queryFn: () => getSprint(sprintId!),
    enabled: Boolean(workspaceId) && Boolean(sprintId),
  })
}
