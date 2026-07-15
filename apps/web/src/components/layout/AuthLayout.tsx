import { Outlet } from 'react-router-dom'

// Thin passthrough: the auth screens render their own full-width `AuthTopBar`
// and the onboarding steps render their own wordmark + progress bar via
// `OnboardingLayout` — this just supplies the shared full-height background.
export function AuthLayout() {
  return (
    <div className="min-h-svh bg-surface-2">
      <Outlet />
    </div>
  )
}
