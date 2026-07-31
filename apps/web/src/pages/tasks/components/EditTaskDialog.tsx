import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import type { Project, Task } from '@foundryhq/shared-types'
import { taskSchema, type TaskFormValues } from '@foundryhq/shared-validation'
import { Button } from '@/components/ui/button'
import { DateField } from '@/components/ui/date-field'
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
import { useDeleteTask } from '@/hooks/useDeleteTask'
import { useSprints } from '@/hooks/useSprints'
import { useUpdateTask } from '@/hooks/useUpdateTask'
import { useWorkspaceMembers } from '@/hooks/useWorkspaceMembers'

const nativeSelectClassName =
  'h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-base outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 md:text-sm dark:bg-input/30'

interface EditTaskDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  task: Task
  projects: Project[]
}

function formValuesFromTask(task: Task): TaskFormValues {
  return {
    title: task.title,
    projectId: task.projectId,
    assigneeId: task.assigneeId ?? '',
    sprintId: task.sprintId ?? '',
    priority: task.priority,
    storyPoints: task.storyPoints != null ? String(task.storyPoints) : '',
    dueDate: task.dueDate ?? '',
  }
}

// Rendered keyed by task.id at the call site (see TasksPage), so this
// component fully remounts — with a fresh form and no confirm-delete state
// left over — every time a different task is opened. Simpler than wiring a
// reset effect for the same result.
export function EditTaskDialog({ open, onOpenChange, task, projects }: EditTaskDialogProps) {
  const updateTask = useUpdateTask()
  const deleteTask = useDeleteTask()
  const { data: members } = useWorkspaceMembers()
  const { data: sprints } = useSprints()
  const [confirmingDelete, setConfirmingDelete] = useState(false)

  const form = useForm({
    resolver: zodResolver(taskSchema),
    defaultValues: formValuesFromTask(task),
  })

  function onSubmit(values: TaskFormValues) {
    updateTask.mutate(
      {
        taskId: task.id,
        title: values.title,
        projectId: values.projectId,
        assigneeId: values.assigneeId || undefined,
        clearAssignee: !values.assigneeId && task.assigneeId != null,
        sprintId: values.sprintId || undefined,
        clearSprint: !values.sprintId && task.sprintId != null,
        priority: values.priority,
        storyPoints: values.storyPoints ? Number(values.storyPoints) : undefined,
        clearStoryPoints: !values.storyPoints && task.storyPoints != null,
        dueDate: values.dueDate || undefined,
        clearDueDate: !values.dueDate && task.dueDate != null,
      },
      { onSuccess: () => onOpenChange(false) }
    )
  }

  function handleDelete() {
    deleteTask.mutate(task.id, { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit task</DialogTitle>
          <DialogDescription>Update this task, or remove it entirely.</DialogDescription>
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
                    <select {...field} className={nativeSelectClassName}>
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
                    <select {...field} className={nativeSelectClassName}>
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
            <div className="grid grid-cols-2 gap-3">
              <FormField
                control={form.control}
                name="priority"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Priority</FormLabel>
                    <FormControl>
                      <select {...field} className={nativeSelectClassName}>
                        <option value="urgent">Urgent</option>
                        <option value="high">High</option>
                        <option value="medium">Medium</option>
                        <option value="low">Low</option>
                      </select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="storyPoints"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Story points</FormLabel>
                    <FormControl>
                      <Input type="number" min={0} step={1} placeholder="Optional" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <FormField
                control={form.control}
                name="sprintId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Sprint</FormLabel>
                    <FormControl>
                      <select {...field} className={nativeSelectClassName}>
                        <option value="">Backlog</option>
                        {sprints?.map((sprint) => (
                          <option key={sprint.id} value={sprint.id}>
                            {sprint.name}
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
                name="dueDate"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Due date</FormLabel>
                    <FormControl>
                      <DateField value={field.value} onChange={field.onChange} placeholder="Optional" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            {updateTask.isError && (
              <p className="text-sm text-destructive">{updateTask.error.message}</p>
            )}
            {deleteTask.isError && (
              <p className="text-sm text-destructive">{deleteTask.error.message}</p>
            )}
            <DialogFooter className="sm:justify-between">
              {confirmingDelete ? (
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">Delete this task?</span>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setConfirmingDelete(false)}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="button"
                    variant="destructive"
                    size="sm"
                    onClick={handleDelete}
                    disabled={deleteTask.isPending}
                  >
                    {deleteTask.isPending ? 'Deleting…' : 'Confirm delete'}
                  </Button>
                </div>
              ) : (
                <Button type="button" variant="destructive" onClick={() => setConfirmingDelete(true)}>
                  Delete
                </Button>
              )}
              <Button type="submit" disabled={updateTask.isPending}>
                {updateTask.isPending ? 'Saving…' : 'Save changes'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
