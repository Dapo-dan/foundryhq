import { ChevronLeft } from 'lucide-react'
import { Navigate, Outlet, useLocation, useNavigate } from 'react-router-dom'
import { Logo } from '@/components/layout/Logo'
import { cn } from '@/lib/utils'
import { useOnboardingStore } from '@/store/slices/onboarding'
import { ONBOARDING_STEPS, type OnboardingStep } from '@/types/onboarding'

// Per design: only workspace/team-size/role/tools/invite are numbered steps
// ("Step N of 5") — sign-up isn't counted, and welcome is an unnumbered
// completion screen with no progress chrome at all.
const NUMBERED_STEPS = ONBOARDING_STEPS.filter((step) => step !== 'welcome')
const TOTAL_STEPS = NUMBERED_STEPS.length

function currentStepFromPath(pathname: string): OnboardingStep | undefined {
  return ONBOARDING_STEPS.find((step) => pathname.endsWith(step))
}

export function OnboardingLayout() {
  const location = useLocation()
  const navigate = useNavigate()
  const completedSteps = useOnboardingStore((state) => state.completedSteps)

  const currentStep = currentStepFromPath(location.pathname)
  if (!currentStep) {
    return <Navigate to="/onboarding/workspace" replace />
  }

  const stepIndex = ONBOARDING_STEPS.indexOf(currentStep)

  // Frontend-only enforcement (no server-side check exists yet): if the user
  // jumps straight to a later step's URL without finishing the ones before
  // it, send them back to the first one they haven't completed.
  const firstIncompleteIndex = ONBOARDING_STEPS.findIndex(
    (step) => !completedSteps.includes(step)
  )
  if (firstIncompleteIndex !== -1 && stepIndex > firstIncompleteIndex) {
    return <Navigate to={`/onboarding/${ONBOARDING_STEPS[firstIncompleteIndex]}`} replace />
  }

  if (currentStep === 'welcome') {
    return (
      <>
        <div className="flex items-center justify-center border-b border-border bg-white px-6 py-4 sm:px-20">
          <Logo />
        </div>
        <div className="flex min-h-[calc(100svh-65px)] items-center justify-center px-4 py-12">
          <div className="flex w-full max-w-[440px] flex-col gap-8">
            <Outlet />
          </div>
        </div>
      </>
    )
  }

  const numberedStepIndex = NUMBERED_STEPS.indexOf(currentStep)
  const stepNumber = numberedStepIndex + 1

  return (
    <>
      <div className="flex items-center justify-between border-b border-border bg-white px-6 py-4 sm:px-20">
        <Logo />
        <span className="text-sm text-text-secondary">
          Step {stepNumber} of {TOTAL_STEPS}
        </span>
      </div>
      <div className="flex min-h-[calc(100svh-65px)] items-center justify-center px-4 py-12">
        <div className="flex w-full max-w-[440px] flex-col gap-8">
          <div className="flex flex-col gap-3">
            <div className="flex gap-1.5">
              {Array.from({ length: TOTAL_STEPS }, (_, i) => (
                <div
                  key={i}
                  className={cn(
                    'h-1 flex-1 rounded-full',
                    i < stepNumber ? 'bg-primary' : 'bg-muted'
                  )}
                />
              ))}
            </div>
            {numberedStepIndex > 0 && (
              <button
                type="button"
                onClick={() => navigate(-1)}
                className="flex items-center gap-1 self-start text-sm text-muted-foreground hover:text-foreground"
              >
                <ChevronLeft className="size-4" />
                Back
              </button>
            )}
          </div>
          <Outlet />
        </div>
      </div>
    </>
  )
}
