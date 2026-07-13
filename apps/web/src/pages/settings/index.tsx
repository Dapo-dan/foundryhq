import { Settings } from 'lucide-react'

export function SettingsPage() {
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-bold">Settings</h1>
        <p className="text-sm text-muted-foreground">
          Manage your workspace, team, and account preferences.
        </p>
      </div>
      <div className="flex flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border py-16 text-center">
        <Settings size={24} className="text-muted-foreground" />
        <h2 className="text-lg font-medium">Nothing to configure yet</h2>
        <p className="max-w-sm text-sm text-muted-foreground">
          Workspace and account settings will appear here.
        </p>
      </div>
    </div>
  )
}
