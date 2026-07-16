import { Building2, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function CompaniesPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">Companies</h1>
          <p className="text-sm text-muted-foreground">
            Manage the companies your team is selling to.
          </p>
        </div>
        <Button>
          <Plus size={20} />
          New company
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <Building2 size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">No companies yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Add a company to start tracking contacts and deals for it.
        </p>
        <Button variant="outline">Add company</Button>
      </div>
    </div>
  )
}
