import { CheckSquare, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function TasksPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Tasks</h1>
          <p className="text-sm text-muted-foreground">
            See and manage the tasks assigned to your workspace.
          </p>
        </div>
        <Button>
          <Plus size={20} />
          New task
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <CheckSquare size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">No tasks yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Create a task to start tracking work for your team.
        </p>
        <Button variant="outline">Create task</Button>
      </div>
    </div>
  )
}
