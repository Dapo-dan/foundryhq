import { useState } from 'react'
import { Zap, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { EmptyState } from '@/components/ui/empty-state'
import { useSprints } from '@/hooks/useSprints'
import { NewSprintDialog } from './components/NewSprintDialog'
import { SprintCard } from './components/SprintCard'

export function SprintsPage() {
  const [dialogOpen, setDialogOpen] = useState(false)
  const { data: sprints, isPending } = useSprints()

  const hasSprints = (sprints?.length ?? 0) > 0

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Sprints</h1>
          <p className="text-sm text-muted-foreground">
            Plan sprints and track your team's velocity.
          </p>
        </div>
        <Button onClick={() => setDialogOpen(true)}>
          <Plus size={20} />
          New sprint
        </Button>
      </div>

      {isPending ? (
        <p className="text-sm text-muted-foreground">Loading sprints…</p>
      ) : hasSprints ? (
        <div className="flex flex-col gap-2">
          {sprints!.map((sprint) => (
            <SprintCard key={sprint.id} sprint={sprint} />
          ))}
        </div>
      ) : (
        <EmptyState
          icon={Zap}
          title="No sprints yet"
          description="Create a sprint to start planning and assigning work."
          action={
            <Button variant="outline" onClick={() => setDialogOpen(true)}>
              Create sprint
            </Button>
          }
        />
      )}

      <NewSprintDialog open={dialogOpen} onOpenChange={setDialogOpen} />
    </div>
  )
}
