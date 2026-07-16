export type TeamSize = 'just_me' | 'small' | 'medium' | 'large'

export type Role = 'founder_ceo' | 'coo_operator' | 'head_of_sales' | 'product_engineer'

// Sign-up (at /auth/sign-up) isn't part of this wizard's numbered progression.
// Of these six, only workspace/team-size/role/tools/invite are numbered steps
// ("Step N of 5") in OnboardingLayout — welcome is an unnumbered completion screen.
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
