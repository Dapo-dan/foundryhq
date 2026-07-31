import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import type { Project } from '@foundryhq/shared-types'
import { taskSchema, type TaskFormValues } from '@foundryhq/shared-validation'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useCreateTask } from '@/hooks/useCreateTask'
import { useWorkspaceMembers } from '@/hooks/useWorkspaceMembers'

interface NewTaskDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  projects: Project[]
}

export function NewTaskDialog({ open, onOpenChange, projects }: NewTaskDialogProps) {
  const createTask = useCreateTask()
  const { data: members } = useWorkspaceMembers()

  const form = useForm({
    resolver: zodResolver(taskSchema),
    defaultValues: { title: '', projectId: '', assigneeId: '' },
  })

  function onSubmit(values: TaskFormValues) {
    createTask.mutate(
      { title: values.title, projectId: values.projectId, assigneeId: values.assigneeId || undefined },
      {
        onSuccess: () => {
          onOpenChange(false)
          form.reset()
        },
      }
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create task</DialogTitle>
          <DialogDescription>Add a task to one of your workspace's projects.</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-3">
            <FormField
              control={form.control}
              name="title"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Title</FormLabel>
                  <FormControl>
                    <Input placeholder="e.g. Write launch announcement" autoFocus {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="projectId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Project</FormLabel>
                  <FormControl>
                    <select
                      {...field}
                      className="h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-base outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 md:text-sm dark:bg-input/30"
                    >
                      <option value="" disabled>
                        Select a project…
                      </option>
                      {projects.map((project) => (
                        <option key={project.id} value={project.id}>
                          {project.name}
                        </option>
                      ))}
                    </select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="assigneeId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Assignee</FormLabel>
                  <FormControl>
                    <select
                      {...field}
                      className="h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-base outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 md:text-sm dark:bg-input/30"
                    >
                      <option value="">Unassigned</option>
                      {members?.map((member) => (
                        <option key={member.id} value={member.userId}>
                          {member.email}
                        </option>
                      ))}
                    </select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {createTask.isError && (
              <p className="text-sm text-destructive">{createTask.error.message}</p>
            )}
            <DialogFooter>
              <Button type="submit" disabled={createTask.isPending}>
                {createTask.isPending ? 'Creating…' : 'Create task'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
