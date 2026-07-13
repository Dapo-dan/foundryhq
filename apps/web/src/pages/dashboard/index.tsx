import { LayoutDashboard, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function DashboardPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Dashboard</h1>
          <p className="text-sm text-muted-foreground">
            An overview of your workspace's activity.
          </p>
        </div>
        <Button>
          <Plus size={20} />
          New
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <LayoutDashboard size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">Nothing to show yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Your workspace overview will appear here once there's activity to report.
        </p>
        <Button variant="outline">Get started</Button>
      </div>
    </div>
  )
}
