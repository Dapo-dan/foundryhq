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
  // Whether this account has ever finished the wizard — distinct from
  // `completedSteps`, which `reset()` clears once the user reaches the
  // dashboard. Sign-in reads this to decide whether to resume onboarding.
  onboardingComplete: boolean
  setWorkspaceName: (name: string) => void
  setTeamSize: (size: TeamSize) => void
  setRole: (role: Role) => void
  toggleTool: (tool: string) => void
  setInvites: (emails: string[]) => void
  markStepComplete: (step: OnboardingStep) => void
  markOnboardingComplete: () => void
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
      onboardingComplete: false,
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
      markOnboardingComplete: () => set({ onboardingComplete: true }),
      // Only clears wizard progress (`initialState`) — `onboardingComplete`
      // is deliberately excluded so it survives this reset.
      reset: () => set(initialState),
    }),
    { name: 'foundryhq-onboarding', storage: createJSONStorage(() => sessionStorage) }
  )
)
