import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  TrendingUp,
  Building2,
  Users,
  SquareCheck,
  Zap,
  ChartBar,
  FileText,
  Settings,
} from 'lucide-react'
import { Logo } from '@/components/layout/Logo'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'

const NAV_SECTIONS = [
  {
    items: [
      { to: '/dashboard', label: 'Overview', icon: LayoutDashboard },
      { to: '/crm/deal-pipeline', label: 'Deal Pipeline', icon: TrendingUp },
      { to: '/crm/companies', label: 'Companies', icon: Building2 },
      { to: '/crm/contacts', label: 'Contacts', icon: Users },
    ],
  },
  {
    title: 'Sprint',
    items: [
      { to: '/tasks', label: 'Tasks', icon: SquareCheck },
      { to: '/sprints', label: 'Sprints', icon: Zap },
    ],
  },
  {
    title: 'Insights',
    items: [
      { to: '/kpis/metrics', label: 'Metrics', icon: ChartBar },
      { to: '/kpis/reports', label: 'Reports', icon: FileText },
    ],
  },
]

export function Sidebar() {
  return (
    <aside className="flex h-svh w-60 shrink-0 flex-col bg-sidebar text-sidebar-foreground">
      <div className="flex h-[60px] shrink-0 items-center border-b border-sidebar-border px-5">
        <Logo />
      </div>

      <nav className="flex flex-1 flex-col gap-0.5 overflow-y-auto p-3">
        {NAV_SECTIONS.map((section, index) => (
          <div key={section.title ?? index} className={index > 0 ? 'mt-2.5' : undefined}>
            {section.title ? (
              <div className="px-3.5 py-1">
                <span className="text-[10px] font-semibold tracking-wide text-muted-foreground uppercase">
                  {section.title}
                </span>
              </div>
            ) : null}
            {section.items.map(({ to, label, icon: Icon }) => (
              <NavLink key={to} to={to}>
                {({ isActive }) => (
                  <div
                    className={cn(
                      'flex h-9 items-center gap-2 rounded-md px-0 text-sm',
                      isActive ? 'bg-brand/10' : 'hover:bg-sidebar-accent',
                    )}
                  >
                    <span
                      className={cn(
                        'h-[18px] w-[3px] shrink-0 rounded-sm',
                        isActive ? 'bg-brand' : 'bg-transparent',
                      )}
                    />
                    <Icon
                      size={15}
                      className={cn('shrink-0', isActive ? 'text-brand' : 'text-muted-foreground')}
                    />
                    <span
                      className={cn(
                        isActive ? 'font-semibold text-foreground' : 'text-muted-foreground',
                      )}
                    >
                      {label}
                    </span>
                  </div>
                )}
              </NavLink>
            ))}
          </div>
        ))}
      </nav>

      <div className="flex shrink-0 items-center gap-2.5 border-t border-sidebar-border px-4 py-4">
        <Avatar>
          <AvatarFallback className="bg-brand text-white">JL</AvatarFallback>
        </Avatar>
        <div className="flex min-w-0 flex-1 flex-col">
          <span className="truncate text-xs font-semibold text-foreground">Jordan Lee</span>
          <span className="truncate text-[11px] text-muted-foreground">jordan@foundryhq.com</span>
        </div>
        <NavLink to="/settings" aria-label="Settings">
          <Settings size={14} className="text-muted-foreground" />
        </NavLink>
      </div>
    </aside>
  )
}
