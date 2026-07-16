import type { DashboardStatus } from '@/types/dashboard'

export const statusTextClass: Record<DashboardStatus, string> = {
  positive: 'text-green-500',
  info: 'text-blue-500',
  warning: 'text-yellow-500',
  danger: 'text-destructive',
}

export const statusBarClass: Record<DashboardStatus, string> = {
  positive: 'bg-green-500',
  info: 'bg-blue-500',
  warning: 'bg-yellow-500',
  danger: 'bg-destructive',
}
