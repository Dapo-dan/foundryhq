import { useQuery } from '@tanstack/react-query'
import { listTasks, type TaskFilters } from '@/services/tasks'
import { useWorkspaceStore } from '@/store/slices/workspace'

export function useTasks(filters: TaskFilters = {}) {
  const workspaceId = useWorkspaceStore((state) => state.currentWorkspaceId)
  return useQuery({
    queryKey: ['tasks', workspaceId, filters],
    queryFn: () => listTasks(filters),
    enabled: Boolean(workspaceId),
  })
}
