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
import { useResetPassword } from '@/hooks/useResetPassword'
import { resetPasswordSchema, type ResetPasswordFormValues } from '@/lib/validation/auth'
import { PasswordStrengthBar } from './PasswordStrengthBar'

export function ResetPasswordForm() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token') ?? ''
  const resetPassword = useResetPassword()
  const form = useForm({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: { password: '', confirmPassword: '' },
  })

  function onSubmit(values: ResetPasswordFormValues) {
    resetPassword.mutate(
      { token, password: values.password },
      { onSuccess: () => navigate('/auth/sign-in') }
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
              <FormLabel>New password</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder="At least 8 characters"
                  autoComplete="new-password"
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

        {resetPassword.isError && (
          <p className="text-sm text-destructive">{resetPassword.error.message}</p>
        )}

        <Button
          type="submit"
          className="mt-1 h-11 w-full text-[15px]"
          disabled={resetPassword.isPending}
        >
          {resetPassword.isPending ? 'Updating…' : 'Update password →'}
        </Button>
      </form>
    </Form>
  )
}
