import AsyncStorage from '@react-native-async-storage/async-storage'
import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'

export type OnboardingStep = 'workspace' | 'invite'

interface OnboardingState {
  workspaceName: string
  invites: string[]
  completedSteps: OnboardingStep[]
  // Whether this account has ever finished onboarding — distinct from
  // `completedSteps`, which only tracks progress through the current run.
  // Sign-in reads this to decide whether to resume onboarding.
  onboardingComplete: boolean
  setWorkspaceName: (name: string) => void
  setInvites: (emails: string[]) => void
  markStepComplete: (step: OnboardingStep) => void
  markOnboardingComplete: () => void
}

// Persisted to AsyncStorage (mobile's analog of web's sessionStorage):
// wizard progress should survive the app being backgrounded or killed
// mid-flow, which happens far more often on mobile than a browser tab
// refreshing.
export const useOnboardingStore = create<OnboardingState>()(
  persist(
    (set) => ({
      workspaceName: '',
      invites: [],
      completedSteps: [],
      onboardingComplete: false,
      setWorkspaceName: (name) => set({ workspaceName: name }),
      setInvites: (emails) => set({ invites: emails }),
      markStepComplete: (step) =>
        set((state) => ({
          completedSteps: state.completedSteps.includes(step)
            ? state.completedSteps
            : [...state.completedSteps, step],
        })),
      markOnboardingComplete: () => set({ onboardingComplete: true }),
    }),
    { name: 'foundryhq-onboarding', storage: createJSONStorage(() => AsyncStorage) }
  )
)
