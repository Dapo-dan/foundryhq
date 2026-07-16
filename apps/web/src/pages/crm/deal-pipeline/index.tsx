import { TrendingUp, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function DealPipelinePage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Deal Pipeline</h1>
          <p className="text-sm text-muted-foreground">
            Track deals as they move from lead to close.
          </p>
        </div>
        <Button>
          <Plus size={20} />
          New deal
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <TrendingUp size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">No deals yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Add a deal to start tracking it through your pipeline stages.
        </p>
        <Button variant="outline">Add deal</Button>
      </div>
    </div>
  )
}
