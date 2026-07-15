import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'
import type { OnboardingStep, Role, TeamSize } from '@/types/onboarding'

interface OnboardingState {
  workspaceName: string
  teamSize: TeamSize | null
  role: Role | null
  tools: string[]
  invites: string[]
  completedSteps: OnboardingStep[]
  setWorkspaceName: (name: string) => void
  setTeamSize: (size: TeamSize) => void
  setRole: (role: Role) => void
  toggleTool: (tool: string) => void
  setInvites: (emails: string[]) => void
  markStepComplete: (step: OnboardingStep) => void
  reset: () => void
}

const initialState = {
  workspaceName: '',
  teamSize: null,
  role: null,
  tools: [],
  invites: [],
  completedSteps: [],
}

// Persisted to sessionStorage (not localStorage): wizard progress should
// survive a mid-flow refresh, but shouldn't outlive the tab or leak across
// different signed-in users on a shared machine.
export const useOnboardingStore = create<OnboardingState>()(
  persist(
    (set) => ({
      ...initialState,
      setWorkspaceName: (name) => set({ workspaceName: name }),
      setTeamSize: (size) => set({ teamSize: size }),
      setRole: (role) => set({ role }),
      toggleTool: (tool) =>
        set((state) => ({
          tools: state.tools.includes(tool)
            ? state.tools.filter((t) => t !== tool)
            : [...state.tools, tool],
        })),
      setInvites: (emails) => set({ invites: emails }),
      markStepComplete: (step) =>
        set((state) => ({
          completedSteps: state.completedSteps.includes(step)
            ? state.completedSteps
            : [...state.completedSteps, step],
        })),
      reset: () => set(initialState),
    }),
    { name: 'foundryhq-onboarding', storage: createJSONStorage(() => sessionStorage) }
  )
)
