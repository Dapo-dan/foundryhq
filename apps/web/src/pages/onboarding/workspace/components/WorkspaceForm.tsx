import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { workspaceSchema, type WorkspaceFormValues } from '@/lib/validation/onboarding'
import { useOnboardingStore } from '@/store/slices/onboarding'

export function WorkspaceForm() {
  const navigate = useNavigate()
  const workspaceName = useOnboardingStore((state) => state.workspaceName)
  const setWorkspaceName = useOnboardingStore((state) => state.setWorkspaceName)
  const markStepComplete = useOnboardingStore((state) => state.markStepComplete)

  const form = useForm({
    resolver: zodResolver(workspaceSchema),
    defaultValues: { name: workspaceName },
  })

  function onSubmit(values: WorkspaceFormValues) {
    setWorkspaceName(values.name)
    markStepComplete('workspace')
    navigate('/onboarding/team-size')
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-3">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input placeholder="e.g. Acme Inc." autoFocus {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" className="mt-1 h-11 w-full text-[15px]">
          Continue →
        </Button>
      </form>
    </Form>
  )
}
