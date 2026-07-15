import { ChevronLeft } from 'lucide-react'
import { Navigate, Outlet, useLocation, useNavigate } from 'react-router-dom'
import { Logo } from '@/components/layout/Logo'
import { cn } from '@/lib/utils'
import { useOnboardingStore } from '@/store/slices/onboarding'
import { ONBOARDING_STEPS, type OnboardingStep } from '@/types/onboarding'

// Sign-up (at /auth/sign-up) is step 1 of 7; the six steps here are 2-7.
const TOTAL_STEPS = ONBOARDING_STEPS.length + 1

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

  const stepNumber = stepIndex + 2 // +1 for 0-index, +1 for sign-up being step 1

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
            {stepIndex > 0 && (
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
