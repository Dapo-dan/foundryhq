import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
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
import { useForgotPassword } from '@/hooks/useForgotPassword'
import { forgotPasswordSchema, type ForgotPasswordFormValues } from '@/lib/validation/auth'

export function ForgotPasswordForm() {
  const [sentTo, setSentTo] = useState<string | null>(null)
  const forgotPassword = useForgotPassword()
  const form = useForm({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: { email: '' },
  })

  function onSubmit(values: ForgotPasswordFormValues) {
    forgotPassword.mutate(values, { onSuccess: () => setSentTo(values.email) })
  }

  if (sentTo) {
    return (
      <p className="text-center text-sm text-muted-foreground">
        If an account exists for <span className="font-medium text-foreground">{sentTo}</span>,
        we've sent a password reset link to that address.
      </p>
    )
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-3">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Your email</FormLabel>
              <FormControl>
                <Input placeholder="you@company.com" autoComplete="email" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {forgotPassword.isError && (
          <p className="text-sm text-destructive">{forgotPassword.error.message}</p>
        )}

        <Button
          type="submit"
          className="mt-1 h-11 w-full text-[15px]"
          disabled={forgotPassword.isPending}
        >
          {forgotPassword.isPending ? 'Sending…' : 'Send reset link →'}
        </Button>
      </form>
    </Form>
  )
}
