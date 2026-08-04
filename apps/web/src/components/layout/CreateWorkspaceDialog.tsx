import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { useQueryClient } from '@tanstack/react-query'
import { workspaceSchema, type WorkspaceFormValues } from '@foundryhq/shared-validation'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useCreateWorkspace } from '@/hooks/useCreateWorkspace'
import { useWorkspaceStore } from '@/store/slices/workspace'

interface CreateWorkspaceDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateWorkspaceDialog({ open, onOpenChange }: CreateWorkspaceDialogProps) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const createWorkspace = useCreateWorkspace()

  const form = useForm({
    resolver: zodResolver(workspaceSchema),
    defaultValues: { name: '' },
  })

  function onSubmit(values: WorkspaceFormValues) {
    createWorkspace.mutate(
      { name: values.name },
      {
        onSuccess: (workspace) => {
          const { workspaces, setWorkspaces, setCurrentWorkspaceId } = useWorkspaceStore.getState()
          setWorkspaces([...workspaces, workspace])
          setCurrentWorkspaceId(workspace.id)
          queryClient.invalidateQueries({ queryKey: ['workspaces'] })
          onOpenChange(false)
          form.reset()
          navigate('/dashboard')
        },
      }
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create workspace</DialogTitle>
          <DialogDescription>Start a new, separate workspace for a different team.</DialogDescription>
        </DialogHeader>
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
            <DialogFooter>
              <Button type="submit" disabled={createWorkspace.isPending}>
                {createWorkspace.isPending ? 'Creating…' : 'Create workspace'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
