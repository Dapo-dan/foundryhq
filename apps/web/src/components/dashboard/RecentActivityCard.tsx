import { Card, CardContent, CardHeader, CardTitle, CardAction } from '@/components/ui/card'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import type { ActivityItem } from '@/types/dashboard'

export function RecentActivityCard({ items }: { items: ActivityItem[] }) {
  return (
    <Card className="h-fit">
      <CardHeader className="border-b">
        <CardTitle>Recent Activity</CardTitle>
        <CardAction>
          <a href="#" className="text-xs text-primary">
            View all →
          </a>
        </CardAction>
      </CardHeader>
      <CardContent className="flex flex-col divide-y px-0">
        {items.map((item) => (
          <div key={item.id} className="flex items-center gap-3 px-4 py-3">
            <Avatar>
              <AvatarFallback
                className="text-white"
                style={{ backgroundColor: item.avatarColor }}
              >
                {item.actorInitials}
              </AvatarFallback>
            </Avatar>
            <div className="flex flex-1 flex-col gap-0.5">
              <span className="text-sm font-semibold">{item.actorName}</span>
              <span className="text-xs text-muted-foreground">{item.action}</span>
            </div>
            <span className="text-xs text-muted-foreground/70">{item.time}</span>
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
