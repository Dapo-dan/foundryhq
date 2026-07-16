import { KpiCard } from '@/components/dashboard/KpiCard'
import { RecentActivityCard } from '@/components/dashboard/RecentActivityCard'
import { GoalsCard } from '@/components/dashboard/GoalsCard'
import { UpcomingCard } from '@/components/dashboard/UpcomingCard'
import {
  mockKpis,
  mockActivity,
  mockGoals,
  mockUpcoming,
  goalsOnTrackSummary,
} from '@/lib/mock/dashboard'

export function DashboardPage() {
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-bold">Overview</h1>
        <p className="text-sm text-muted-foreground">
          An overview of your workspace's activity.
        </p>
      </div>

      <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {mockKpis.map((metric) => (
          <KpiCard key={metric.id} metric={metric} />
        ))}
      </div>

      <div className="grid grid-cols-1 items-start gap-5 lg:grid-cols-[1fr_320px]">
        <RecentActivityCard items={mockActivity} />
        <div className="flex flex-col gap-4">
          <GoalsCard goals={mockGoals} summary={goalsOnTrackSummary} />
          <UpcomingCard events={mockUpcoming} />
        </div>
      </div>
    </div>
  )
}
