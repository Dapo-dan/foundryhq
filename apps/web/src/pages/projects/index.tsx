import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { FolderKanban, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { EmptyState } from '@/components/ui/empty-state'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useCreateProject } from '@/hooks/useCreateProject'
import { useProjects } from '@/hooks/useProjects'
import { projectSchema, type ProjectFormValues } from '@foundryhq/shared-validation'
import { ProjectCard } from './components/ProjectCard'

export function ProjectsPage() {
  const [open, setOpen] = useState(false)
  const { data: projects, isPending } = useProjects()
  const createProject = useCreateProject()

  const form = useForm({
    resolver: zodResolver(projectSchema),
    defaultValues: { name: '', description: '' },
  })

  function onSubmit(values: ProjectFormValues) {
    createProject.mutate(values, {
      onSuccess: () => {
        setOpen(false)
        form.reset()
      },
    })
  }

  const hasProjects = (projects?.length ?? 0) > 0

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Projects</h1>
          <p className="text-sm text-muted-foreground">
            Track the projects your team is working on.
          </p>
        </div>
        <Button onClick={() => setOpen(true)}>
          <Plus size={20} />
          New project
        </Button>
      </div>

      {isPending ? (
        <p className="text-sm text-muted-foreground">Loading projects…</p>
      ) : hasProjects ? (
        <div className="flex flex-col gap-2">
          {projects!.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      ) : (
        <EmptyState
          icon={FolderKanban}
          title="No projects yet"
          description="Create your first project to start organizing work for your team."
          action={
            <Button variant="outline" onClick={() => setOpen(true)}>
              Create project
            </Button>
          }
        />
      )}

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create project</DialogTitle>
            <DialogDescription>Give your project a name to get started.</DialogDescription>
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
                      <Input placeholder="e.g. Website Redesign" autoFocus {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <Input placeholder="Optional" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              {createProject.isError && (
                <p className="text-sm text-destructive">{createProject.error.message}</p>
              )}
              <DialogFooter>
                <Button type="submit" disabled={createProject.isPending}>
                  {createProject.isPending ? 'Creating…' : 'Create project'}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
