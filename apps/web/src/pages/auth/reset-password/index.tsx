import { AuthCard } from '@/components/layout/AuthCard'
import { AuthTopBar } from '@/components/layout/AuthTopBar'
import { ResetPasswordForm } from './components/ResetPasswordForm'

export function ResetPasswordPage() {
  return (
    <>
      <AuthTopBar navLabel="← Back to sign in" navHref="/auth/sign-in" />
      <div className="flex min-h-[calc(100svh-65px)] items-center justify-center px-4 py-12">
        <div className="w-full max-w-[440px]">
          <AuthCard heading="Set a new password" description="Choose a new password for your account.">
            <ResetPasswordForm />
          </AuthCard>
        </div>
      </div>
    </>
  )
}
