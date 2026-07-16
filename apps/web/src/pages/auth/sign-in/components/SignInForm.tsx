import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/ui/password-input'
import { Separator } from '@/components/ui/separator'
import { useSignIn } from '@/hooks/useSignIn'
import { signInSchema, type SignInFormValues } from '@/lib/validation/auth'
import { useOnboardingStore } from '@/store/slices/onboarding'
import { ONBOARDING_STEPS } from '@/types/onboarding'
import { OAuthButtons } from '../../components/OAuthButtons'

export function SignInForm() {
  const navigate = useNavigate()
  const signIn = useSignIn()
  const form = useForm({
    resolver: zodResolver(signInSchema),
    defaultValues: { email: '', password: '' },
  })

  function onSubmit(values: SignInFormValues) {
    signIn.mutate(values, {
      onSuccess: () => {
        const { onboardingComplete, completedSteps } = useOnboardingStore.getState()
        if (onboardingComplete) {
          navigate('/dashboard')
          return
        }
        // Resume exactly where they left off, same rule OnboardingLayout's
        // step-guard uses to find the first step they haven't finished yet.
        const nextStep = ONBOARDING_STEPS.find((step) => !completedSteps.includes(step))
        navigate(`/onboarding/${nextStep ?? 'workspace'}`)
      },
    })
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-3">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Work email</FormLabel>
              <FormControl>
                <Input placeholder="you@company.com" autoComplete="email" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <div className="flex items-center justify-between">
                <FormLabel>Password</FormLabel>
                <Link
                  to="/auth/forgot-password"
                  className="text-xs text-brand hover:text-brand-accent"
                >
                  Forgot password?
                </Link>
              </div>
              <FormControl>
                <PasswordInput
                  placeholder="At least 8 characters"
                  autoComplete="current-password"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {signIn.isError && <p className="text-sm text-destructive">{signIn.error.message}</p>}

        <Button type="submit" className="mt-1 h-11 w-full text-[15px]" disabled={signIn.isPending}>
          {signIn.isPending ? 'Signing in…' : 'Sign in →'}
        </Button>

        <div className="my-1 flex items-center gap-2.5">
          <Separator className="flex-1" />
          <span className="text-xs text-text-subtle">or continue with</span>
          <Separator className="flex-1" />
        </div>

        <OAuthButtons />
      </form>
    </Form>
  )
}
