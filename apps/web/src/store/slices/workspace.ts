import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'
import type { Workspace } from '@foundryhq/shared-types'

interface WorkspaceState {
  workspaces: Workspace[]
  currentWorkspaceId: string | null
  setWorkspaces: (workspaces: Workspace[]) => void
  setCurrentWorkspaceId: (id: string) => void
  clear: () => void
}

// Only currentWorkspaceId is persisted (to localStorage, so the choice
// survives across sessions) — workspaces itself always comes fresh from
// GET /workspaces via useWorkspaces, never from storage.
export const useWorkspaceStore = create<WorkspaceState>()(
  persist(
    (set) => ({
      workspaces: [],
      currentWorkspaceId: null,
      setWorkspaces: (workspaces) =>
        set((state) => ({
          workspaces,
          currentWorkspaceId: workspaces.some((w) => w.id === state.currentWorkspaceId)
            ? state.currentWorkspaceId
            : workspaces[0]?.id ?? null,
        })),
      setCurrentWorkspaceId: (id) => set({ currentWorkspaceId: id }),
      clear: () => set({ workspaces: [], currentWorkspaceId: null }),
    }),
    {
      name: 'foundryhq-workspace',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({ currentWorkspaceId: state.currentWorkspaceId }),
    }
  )
)
