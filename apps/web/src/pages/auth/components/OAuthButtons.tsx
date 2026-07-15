import { Button } from '@/components/ui/button'

// Disabled placeholders: OAuth is explicitly deferred (see the commented-out
// google/github vars in apps/api/.env.example) — no provider is wired up yet.
// The mockup uses plain glyphs here rather than brand logos (lucide-react
// dropped brand/logo icons), so this matches that rather than reaching for a
// new icon dependency for two disabled buttons.
export function OAuthButtons() {
  return (
    <div className="grid grid-cols-2 gap-3">
      <Button type="button" variant="outline" className="h-11 gap-2" disabled title="Coming soon">
        <span aria-hidden className="text-sm font-semibold">G</span>
        Google
      </Button>
      <Button type="button" variant="outline" className="h-11 gap-2" disabled title="Coming soon">
        <span aria-hidden className="text-sm">◆</span>
        GitHub
      </Button>
    </div>
  )
}
