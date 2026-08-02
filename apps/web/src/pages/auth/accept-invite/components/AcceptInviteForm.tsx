import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { PasswordInput } from '@/components/ui/password-input'
import { useAcceptInvite } from '@/hooks/useAcceptInvite'
import { acceptInviteSchema, type AcceptInviteFormValues } from '@foundryhq/shared-validation'
import { PasswordStrengthBar } from '../../reset-password/components/PasswordStrengthBar'

export function AcceptInviteForm() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token') ?? ''
  const acceptInvite = useAcceptInvite()
  const form = useForm({
    resolver: zodResolver(acceptInviteSchema),
    defaultValues: { password: '', confirmPassword: '' },
  })

  function onSubmit(values: AcceptInviteFormValues) {
    acceptInvite.mutate(
      { token, password: values.password },
      // Already logged in via the token — go straight to the workspace
      // instead of back through sign-in, unlike ResetPasswordForm.
      { onSuccess: () => navigate('/dashboard') }
    )
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-3">
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
                  autoFocus
                  {...field}
                />
              </FormControl>
              <PasswordStrengthBar password={field.value} />
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="confirmPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Confirm password</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder="Re-enter your password"
                  autoComplete="new-password"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {acceptInvite.isError && (
          <p className="text-sm text-destructive">{acceptInvite.error.message}</p>
        )}

        <Button
          type="submit"
          className="mt-1 h-11 w-full text-[15px]"
          disabled={acceptInvite.isPending}
        >
          {acceptInvite.isPending ? 'Setting up your account…' : 'Set password & continue →'}
        </Button>
      </form>
    </Form>
  )
}
