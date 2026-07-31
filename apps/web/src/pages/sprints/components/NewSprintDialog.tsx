import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { sprintSchema, type SprintFormValues } from '@foundryhq/shared-validation'
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
import { useCreateSprint } from '@/hooks/useCreateSprint'

interface NewSprintDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function NewSprintDialog({ open, onOpenChange }: NewSprintDialogProps) {
  const createSprint = useCreateSprint()

  const form = useForm({
    resolver: zodResolver(sprintSchema),
    defaultValues: { name: '', startDate: '', endDate: '' },
  })

  function onSubmit(values: SprintFormValues) {
    createSprint.mutate(values, {
      onSuccess: () => {
        onOpenChange(false)
        form.reset()
      },
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create sprint</DialogTitle>
          <DialogDescription>Plan a time-boxed sprint for your team's tasks.</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-3">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="e.g. Sprint 12" autoFocus {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="grid grid-cols-2 gap-3">
              <FormField
                control={form.control}
                name="startDate"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Start date</FormLabel>
                    <FormControl>
                      <DateField value={field.value} onChange={field.onChange} placeholder="Start date" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="endDate"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>End date</FormLabel>
                    <FormControl>
                      <DateField value={field.value} onChange={field.onChange} placeholder="End date" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            {createSprint.isError && (
              <p className="text-sm text-destructive">{createSprint.error.message}</p>
            )}
            <DialogFooter>
              <Button type="submit" disabled={createSprint.isPending}>
                {createSprint.isPending ? 'Creating…' : 'Create sprint'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
