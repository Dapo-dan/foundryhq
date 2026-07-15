export type TeamSize = 'just_me' | 'small' | 'medium' | 'large'

export type Role = 'founder_ceo' | 'coo_operator' | 'head_of_sales' | 'product_engineer'

// Sign-up is step 1 of 7 but lives at /auth/sign-up, not under OnboardingLayout —
// these are the six wizard steps that follow it.
export type OnboardingStep = 'workspace' | 'team-size' | 'role' | 'tools' | 'invite' | 'welcome'

export const ONBOARDING_STEPS: OnboardingStep[] = [
  'workspace',
  'team-size',
  'role',
  'tools',
  'invite',
  'welcome',
]

export interface Workspace {
  id: string
  name: string
}

export interface CreateWorkspaceInput {
  name: string
  teamSize: TeamSize
  role: Role
  tools: string[]
}
