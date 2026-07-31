import { useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { listWorkspaces } from '@/services/workspace'
import { useWorkspaceStore } from '@/store/slices/workspace'

// Populates store/slices/workspace.ts as a side effect so components that
// need the current workspace (e.g. TopNav, services/projects.ts) don't each
// need their own query — every AppLayout render (post-login/bootstrap)
// refreshes it here.
export function useWorkspaces(enabled: boolean) {
  const query = useQuery({ queryKey: ['workspaces'], queryFn: listWorkspaces, enabled })

  useEffect(() => {
    if (query.data) {
      useWorkspaceStore.getState().setWorkspaces(query.data)
    }
  }, [query.data])

  return query
}
