import type { LucideIcon } from 'lucide-react'

export type DashboardStatus = 'positive' | 'info' | 'warning' | 'danger'

export interface KpiMetric {
  id: string
  label: string
  icon: LucideIcon
  value: string
  change: string
  status: DashboardStatus
}

export interface ActivityItem {
  id: string
  actorName: string
  actorInitials: string
  avatarColor: string
  action: string
  time: string
}

export interface GoalProgress {
  id: string
  name: string
  percent: number
  status: DashboardStatus
}

export interface UpcomingEvent {
  id: string
  name: string
  time: string
  status: DashboardStatus
}
