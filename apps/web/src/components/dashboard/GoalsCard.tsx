import { Card, CardContent, CardHeader, CardTitle, CardAction } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { statusBarClass, statusTextClass } from '@/components/dashboard/status'
import type { GoalProgress } from '@/types/dashboard'

export function GoalsCard({ goals, summary }: { goals: GoalProgress[]; summary: string }) {
  return (
    <Card>
      <CardHeader className="border-b">
        <CardTitle>Goals</CardTitle>
        <CardAction>
          <span className="text-xs font-medium text-green-500">{summary}</span>
        </CardAction>
      </CardHeader>
      <CardContent className="flex flex-col divide-y px-0">
        {goals.map((goal) => (
          <div key={goal.id} className="flex flex-col gap-2 px-4 py-3">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium">{goal.name}</span>
              <span className={`text-xs font-semibold ${statusTextClass[goal.status]}`}>
                {goal.percent}%
              </span>
            </div>
            <Progress value={goal.percent} indicatorClassName={statusBarClass[goal.status]} />
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
