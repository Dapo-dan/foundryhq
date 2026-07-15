import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
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
import { useSignUp } from '@/hooks/useSignUp'
import { signUpSchema, type SignUpFormValues } from '@/lib/validation/auth'
import { OAuthButtons } from '../../components/OAuthButtons'

export function SignUpForm() {
  const navigate = useNavigate()
  const signUp = useSignUp()
  const form = useForm({
    resolver: zodResolver(signUpSchema),
    defaultValues: { email: '', password: '' },
  })

  function onSubmit(values: SignUpFormValues) {
    signUp.mutate(values, { onSuccess: () => navigate('/onboarding/workspace') })
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
              <FormLabel>Password</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder="At least 8 characters"
                  autoComplete="new-password"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {signUp.isError && <p className="text-sm text-destructive">{signUp.error.message}</p>}

        <Button type="submit" className="mt-1 h-11 w-full text-[15px]" disabled={signUp.isPending}>
          {signUp.isPending ? 'Creating account…' : 'Create account →'}
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
