import { ChartBar, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function MetricsPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Metrics</h1>
          <p className="text-sm text-muted-foreground">
            Define KPIs and track their progress over time.
          </p>
        </div>
        <Button>
          <Plus size={20} />
          New KPI
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <ChartBar size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">No KPIs yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Define a KPI and start recording values against its target.
        </p>
        <Button variant="outline">Add KPI</Button>
      </div>
    </div>
  )
}
