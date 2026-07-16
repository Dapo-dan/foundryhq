import { CheckCircle2 } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { useOnboardingStore } from '@/store/slices/onboarding'
import type { Role, TeamSize } from '@/types/onboarding'

const TEAM_SIZE_LABELS: Record<TeamSize, string> = {
  just_me: 'Just me',
  small: '2-10 people',
  medium: '11-50 people',
  large: '50+ people',
}

const ROLE_LABELS: Record<Role, string> = {
  founder_ceo: 'Founder / CEO',
  coo_operator: 'COO / Operator',
  head_of_sales: 'Head of Sales',
  product_engineer: 'Product / Engineer',
}

export function WelcomeStepPage() {
  const navigate = useNavigate()
  const { workspaceName, teamSize, role, invites, reset } = useOnboardingStore((state) => state)

  const checklist = [
    `"${workspaceName}" workspace created`,
    [teamSize && TEAM_SIZE_LABELS[teamSize], role && ROLE_LABELS[role]].filter(Boolean).join(' · '),
    invites.length > 0
      ? `${invites.length} teammate${invites.length === 1 ? '' : 's'} invited`
      : 'You can invite teammates anytime from Settings',
  ]

  function onGoToDashboard() {
    // Navigate first, then clear the wizard state on the next tick — resetting
    // first re-renders OnboardingLayout (still mounted on this route) with an
    // empty `completedSteps`, and its step-guard redirects back to the first
    // step before the navigation to /dashboard has a chance to take effect.
    navigate('/dashboard')
    setTimeout(reset, 0)
  }

  return (
    <div className="flex flex-col items-center gap-6 text-center">
      <div className="flex size-12 items-center justify-center rounded-full bg-primary">
        <CheckCircle2 className="size-6 text-primary-foreground" />
      </div>
      <div className="flex flex-col gap-1">
        <h1 className="text-2xl font-bold text-foreground">Your workspace is ready!</h1>
        <p className="text-sm text-muted-foreground">
          Everything's set up — here's a quick summary before you dive in.
        </p>
      </div>
      <ul className="flex w-full flex-col gap-2 text-left">
        {checklist.map((item) => (
          <li
            key={item}
            className="flex items-center gap-2 rounded-lg border border-input px-4 py-3 text-sm text-foreground"
          >
            <CheckCircle2 className="size-4 shrink-0 text-brand-accent" />
            {item}
          </li>
        ))}
      </ul>
      <Button type="button" className="h-11 w-full text-[15px]" onClick={onGoToDashboard}>
        Go to dashboard →
      </Button>
    </div>
  )
}
