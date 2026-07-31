import { format } from 'date-fns'
import { Link } from 'react-router-dom'
import type { Sprint } from '@foundryhq/shared-types'
import { useVelocity } from '@/hooks/useVelocity'

function formatDate(dateOnly: string): string {
  return format(new Date(`${dateOnly}T00:00:00`), 'MMM d, yyyy')
}

interface SprintCardProps {
  sprint: Sprint
}

export function SprintCard({ sprint }: SprintCardProps) {
  const { data: velocity } = useVelocity(sprint.id)

  return (
    <Link
      to={`/sprints/${sprint.id}`}
      className="flex items-center justify-between rounded-lg border border-border p-4 hover:bg-muted/30"
    >
      <div>
        <h3 className="font-medium">{sprint.name}</h3>
        <p className="text-sm text-muted-foreground">
          {formatDate(sprint.startDate)} – {formatDate(sprint.endDate)}
        </p>
      </div>
      <div className="text-right">
        <p className="text-lg font-semibold">{(velocity?.velocity ?? 0).toFixed(1)}</p>
        <p className="text-xs text-muted-foreground">velocity</p>
      </div>
    </Link>
  )
}
