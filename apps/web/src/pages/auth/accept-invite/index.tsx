import { AuthCard } from '@/components/layout/AuthCard'
import { AuthTopBar } from '@/components/layout/AuthTopBar'
import { AcceptInviteForm } from './components/AcceptInviteForm'

export function AcceptInvitePage() {
  return (
    <>
      <AuthTopBar navLabel="← Back to sign in" navHref="/auth/sign-in" />
      <div className="flex min-h-[calc(100svh-65px)] items-center justify-center px-4 py-12">
        <div className="w-full max-w-[440px]">
          <AuthCard
            heading="You've been invited"
            description="Set a password to activate your account and join the workspace."
          >
            <AcceptInviteForm />
          </AuthCard>
        </div>
      </div>
    </>
  )
}
