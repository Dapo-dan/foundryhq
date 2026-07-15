import { Link } from 'react-router-dom'
import { AuthCard } from '@/components/layout/AuthCard'
import { AuthTopBar } from '@/components/layout/AuthTopBar'
import { SignUpForm } from './components/SignUpForm'

export function SignUpPage() {
  return (
    <>
      <AuthTopBar navLabel="Already have an account? Sign in →" navHref="/auth/sign-in" />
      <div className="flex min-h-[calc(100svh-65px)] flex-col items-center justify-center gap-4 px-4 py-12">
        <div className="w-full max-w-[440px]">
          <AuthCard
            heading="Create your account"
            description="14-day free trial. No credit card required."
          >
            <SignUpForm />
          </AuthCard>
        </div>
        <p className="max-w-[420px] text-center text-xs text-text-subtle">
          By signing up you agree to our{' '}
          <Link to="/terms" className="text-brand hover:text-brand-accent">
            Terms of Service
          </Link>{' '}
          and{' '}
          <Link to="/privacy" className="text-brand hover:text-brand-accent">
            Privacy Policy
          </Link>
        </p>
      </div>
    </>
  )
}
