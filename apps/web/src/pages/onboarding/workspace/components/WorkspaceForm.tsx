import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { workspaceSchema, type WorkspaceFormValues } from '@foundryhq/shared-validation'
import { useCreateWorkspace } from '@/hooks/useCreateWorkspace'
import { useOnboardingStore } from '@/store/slices/onboarding'

export function WorkspaceForm() {
  const navigate = useNavigate()
  const workspaceName = useOnboardingStore((state) => state.workspaceName)
  const setWorkspaceName = useOnboardingStore((state) => state.setWorkspaceName)
  const setWorkspaceId = useOnboardingStore((state) => state.setWorkspaceId)
  const markStepComplete = useOnboardingStore((state) => state.markStepComplete)
  const createWorkspace = useCreateWorkspace()

  const form = useForm({
    resolver: zodResolver(workspaceSchema),
    defaultValues: { name: workspaceName },
  })

  function onSubmit(values: WorkspaceFormValues) {
    createWorkspace.mutate(
      { name: values.name },
      {
        onSuccess: (workspace) => {
          setWorkspaceName(values.name)
          setWorkspaceId(workspace.id)
          markStepComplete('workspace')
          navigate('/onboarding/team-size')
        },
      }
    )
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
        {createWorkspace.isError && (
          <p className="text-sm text-destructive">{createWorkspace.error.message}</p>
        )}
        <Button
          type="submit"
          className="mt-1 h-11 w-full text-[15px]"
          disabled={createWorkspace.isPending}
        >
          {createWorkspace.isPending ? 'Creating…' : 'Continue →'}
        </Button>
      </form>
    </Form>
  )
}
