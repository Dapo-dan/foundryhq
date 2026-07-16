import { Zap, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function SprintsPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Sprints</h1>
          <p className="text-sm text-muted-foreground">
            Plan sprints and track your team's velocity.
          </p>
        </div>
        <Button>
          <Plus size={20} />
          New sprint
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <Zap size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">No sprints yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Create a sprint to start planning and assigning work.
        </p>
        <Button variant="outline">Create sprint</Button>
      </div>
    </div>
  )
}
