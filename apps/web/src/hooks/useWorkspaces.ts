import { useQuery } from '@tanstack/react-query'
import { listWorkspaces } from '@/services/workspace'

export function useWorkspaces(enabled: boolean) {
  return useQuery({ queryKey: ['workspaces'], queryFn: listWorkspaces, enabled })
}
