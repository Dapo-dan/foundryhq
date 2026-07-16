import { Card, CardContent } from '@/components/ui/card'
import { statusTextClass } from '@/components/dashboard/status'
import type { KpiMetric } from '@/types/dashboard'

export function KpiCard({ metric }: { metric: KpiMetric }) {
  const Icon = metric.icon

  return (
    <Card>
      <CardContent className="flex flex-col gap-2">
        <div className="flex items-center justify-between">
          <span className="text-xs font-medium text-muted-foreground">{metric.label}</span>
          <Icon size={16} className="text-muted-foreground/60" />
        </div>
        <span className="text-2xl font-bold">{metric.value}</span>
        <span className={`text-xs ${statusTextClass[metric.status]}`}>{metric.change}</span>
      </CardContent>
    </Card>
  )
}
