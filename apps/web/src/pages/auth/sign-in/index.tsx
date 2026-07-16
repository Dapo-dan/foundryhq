import { AuthCard } from '@/components/layout/AuthCard'
import { AuthTopBar } from '@/components/layout/AuthTopBar'
import { SignInForm } from './components/SignInForm'

export function SignInPage() {
  return (
    <>
      <AuthTopBar
        navLabel="Don't have an account? Sign up →"
        navLabelCompact="Sign up →"
        navHref="/auth/sign-up"
      />
      <div className="flex min-h-[calc(100svh-65px)] items-center justify-center px-4 py-12">
        <div className="w-full max-w-[440px]">
          <AuthCard heading="Welcome back" description="Sign in to your FoundryHQ workspace.">
            <SignInForm />
          </AuthCard>
        </div>
      </div>
    </>
  )
}
