import { Briefcase, SquareCheck, Target, TrendingUp } from 'lucide-react'
import type { ActivityItem, GoalProgress, KpiMetric, UpcomingEvent } from '@/types/dashboard'

export const mockKpis: KpiMetric[] = [
  {
    id: 'revenue',
    label: 'Monthly Revenue',
    icon: TrendingUp,
    value: '$42.5k',
    change: '+12% vs last month',
    status: 'positive',
  },
  {
    id: 'pipeline',
    label: 'Active Pipeline',
    icon: Briefcase,
    value: '$537k',
    change: '18 open deals',
    status: 'info',
  },
  {
    id: 'tasks',
    label: 'Tasks Due Today',
    icon: SquareCheck,
    value: '7',
    change: '3 overdue',
    status: 'warning',
  },
  {
    id: 'goals',
    label: 'Goals On Track',
    icon: Target,
    value: '4 / 6',
    change: '2 at risk',
    status: 'danger',
  },
]

export const mockActivity: ActivityItem[] = [
  {
    id: '1',
    actorName: 'Sofia Torres',
    actorInitials: 'ST',
    avatarColor: '#581C87',
    action: 'moved Stripe to Proposal',
    time: '2h ago',
  },
  {
    id: '2',
    actorName: 'Jake Nwosu',
    actorInitials: 'JN',
    avatarColor: '#7C2D12',
    action: 'completed API auth flow',
    time: '3h ago',
  },
  {
    id: '3',
    actorName: 'Marcus Reid',
    actorInitials: 'MR',
    avatarColor: '#064E3B',
    action: 'updated KPI · MRR to $42.5k',
    time: '5h ago',
  },
  {
    id: '4',
    actorName: 'Amara Chen',
    actorInitials: 'AC',
    avatarColor: '#1E3A5F',
    action: 'scheduled Q2 Review meeting',
    time: 'Yesterday',
  },
  {
    id: '5',
    actorName: 'Jake Nwosu',
    actorInitials: 'JN',
    avatarColor: '#7C2D12',
    action: 'created Sprint 14',
    time: 'Yesterday',
  },
]

export const mockGoals: GoalProgress[] = [
  { id: '1', name: 'Close $500k ARR', percent: 85, status: 'positive' },
  { id: '2', name: 'Launch v2.0', percent: 62, status: 'info' },
  { id: '3', name: 'Hire 3 Engineers', percent: 33, status: 'warning' },
  { id: '4', name: 'Reach 1k MAU', percent: 18, status: 'danger' },
]

export const mockUpcoming: UpcomingEvent[] = [
  { id: '1', name: 'Q2 Review', time: 'Today · 3:00 PM', status: 'info' },
  { id: '2', name: 'Sprint 14 Kickoff', time: 'Tomorrow · 10:00 AM', status: 'positive' },
  { id: '3', name: 'Investor Update', time: 'Thu · 2:00 PM', status: 'warning' },
]

export const goalsOnTrackSummary = '4 / 6 on track'
