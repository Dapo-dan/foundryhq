import { AuthCard } from '@/components/layout/AuthCard'
import { AuthTopBar } from '@/components/layout/AuthTopBar'
import { ForgotPasswordForm } from './components/ForgotPasswordForm'

export function ForgotPasswordPage() {
  return (
    <>
      <AuthTopBar navLabel="← Back to sign in" navHref="/auth/sign-in" />
      <div className="flex min-h-[calc(100svh-65px)] items-center justify-center px-4 py-12">
        <div className="w-full max-w-[440px]">
          <AuthCard
            heading="Reset your password"
            description="Enter your email address and we'll send you a reset link."
          >
            <ForgotPasswordForm />
          </AuthCard>
        </div>
      </div>
    </>
  )
}
