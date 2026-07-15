import { WorkspaceForm } from './components/WorkspaceForm'

export function WorkspaceStepPage() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1 text-center">
        <h1 className="text-2xl font-bold text-foreground">Name your workspace</h1>
        <p className="text-sm text-muted-foreground">You can always change this later.</p>
      </div>
      <WorkspaceForm />
    </div>
  )
}
