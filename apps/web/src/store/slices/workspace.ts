import { create } from 'zustand'
import type { Workspace } from '@foundryhq/shared-types'

interface WorkspaceState {
  workspaces: Workspace[]
  // Multi-workspace switching isn't built yet — there's always exactly one
  // per user (Phase 1's onboarding always creates one), so this just tracks
  // whichever one setWorkspaces last saw.
  currentWorkspaceId: string | null
  setWorkspaces: (workspaces: Workspace[]) => void
  clear: () => void
}

export const useWorkspaceStore = create<WorkspaceState>((set) => ({
  workspaces: [],
  currentWorkspaceId: null,
  setWorkspaces: (workspaces) =>
    set((state) => ({
      workspaces,
      currentWorkspaceId: workspaces.some((w) => w.id === state.currentWorkspaceId)
        ? state.currentWorkspaceId
        : workspaces[0]?.id ?? null,
    })),
  clear: () => set({ workspaces: [], currentWorkspaceId: null }),
}))
