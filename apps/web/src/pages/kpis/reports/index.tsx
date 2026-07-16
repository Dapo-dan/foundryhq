import { FileText } from 'lucide-react'

export function ReportsPage() {
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-bold">Reports</h1>
        <p className="text-sm text-muted-foreground">
          Shareable snapshots of your workspace's KPIs and goals.
        </p>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <FileText size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">No reports yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Reports will appear here once your KPIs have data to summarize.
        </p>
      </div>
    </div>
  )
}
