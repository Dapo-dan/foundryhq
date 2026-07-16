import { Card, CardContent, CardHeader, CardTitle, CardAction } from '@/components/ui/card'
import { statusBarClass } from '@/components/dashboard/status'
import type { UpcomingEvent } from '@/types/dashboard'

export function UpcomingCard({ events }: { events: UpcomingEvent[] }) {
  return (
    <Card>
      <CardHeader className="border-b">
        <CardTitle>Upcoming</CardTitle>
        <CardAction>
          <a href="#" className="text-xs text-primary">
            View all →
          </a>
        </CardAction>
      </CardHeader>
      <CardContent className="flex flex-col divide-y px-0">
        {events.map((event) => (
          <div key={event.id} className="flex items-center gap-2.5 px-4 py-3">
            <span className={`size-2 shrink-0 rounded-full ${statusBarClass[event.status]}`} />
            <div className="flex flex-col gap-0.5">
              <span className="text-xs font-semibold">{event.name}</span>
              <span className="text-[11px] text-muted-foreground/70">{event.time}</span>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
