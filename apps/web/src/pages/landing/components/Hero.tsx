import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'

function DealCard() {
  return (
    <div className="w-[240px] rounded-xl border border-border bg-white p-4 text-left shadow-[0_8px_24px_-4px_rgba(10,10,10,0.08)]">
      <div className="flex items-center justify-between">
        <span className="text-[13px] font-semibold text-text-primary">Acme Corp</span>
        <span className="rounded-md bg-[#EFF6FF] px-2 py-0.5 text-[10px] font-medium text-brand-accent">
          Proposal
        </span>
      </div>
      <p className="mt-2.5 text-2xl font-extrabold text-text-primary">$48,000</p>
      <div className="mt-1.5 flex items-center gap-1.5">
        <span className="h-1.5 w-1.5 rounded-full bg-[#22C55E]" />
        <span className="text-[11px] text-text-muted">Active · closes Dec 15</span>
      </div>
    </div>
  )
}

function KpiCard() {
  return (
    <div className="w-[200px] rounded-xl border border-border bg-white p-4 text-left shadow-[0_8px_24px_-4px_rgba(10,10,10,0.08)]">
      <p className="text-[11px] leading-snug text-text-muted">Monthly Recurring Revenue</p>
      <p className="mt-1 text-[28px] font-extrabold text-text-primary">$12.4K</p>
      <div className="mt-1.5 flex items-center gap-1.5">
        <span className="rounded-md bg-[#F0FDF4] px-2 py-0.5 text-[11px] font-semibold text-[#16A34A]">
          ↑ 23%
        </span>
        <span className="text-[11px] text-text-muted">vs last month</span>
      </div>
    </div>
  )
}

function SprintCard() {
  return (
    <div className="w-[240px] rounded-xl border border-border bg-white p-4 text-left shadow-[0_8px_24px_-4px_rgba(10,10,10,0.08)]">
      <div className="flex items-center justify-between">
        <span className="text-[13px] font-semibold text-text-primary">Auth flow</span>
        <span className="rounded-md bg-[#F0FDF4] px-2 py-0.5 text-[10px] font-semibold text-[#16A34A]">
          ✓ Done
        </span>
      </div>
      <div className="mt-2.5 flex flex-col gap-0.5">
        <span className="text-[11px] text-text-muted">Assigned to Amara Chen</span>
        <span className="text-[11px] text-text-muted">Sprint 4 · Q1 2026</span>
      </div>
    </div>
  )
}

export function Hero() {
  return (
    <section className="flex flex-col items-center gap-7 px-6 pb-20 pt-16 text-center sm:px-10 lg:px-[100px] lg:pt-24">
      <div className="flex items-center gap-2 rounded-full bg-[#EFF6FF] px-3.5 py-1.5">
        <span className="h-1.5 w-1.5 rounded-full bg-brand" />
        <span className="text-[13px] font-medium text-brand-accent">Built for founding teams</span>
      </div>
      <h1 className="max-w-3xl text-4xl font-extrabold leading-[1.05] text-text-primary sm:text-5xl lg:text-[64px]">
        One workspace.
        <br />
        Every startup need.
      </h1>
      <p className="max-w-xl text-base leading-relaxed text-text-secondary sm:text-lg">
        Replace your CRM, task tracker, meeting notes, and KPI dashboard — one fast workspace built
        for founding teams, from day zero to Series B.
      </p>
      <div className="flex flex-wrap items-center justify-center gap-3.5">
        <Button asChild size="lg" className="px-8 text-[15px]">
          <Link to="/onboarding">Start for free →</Link>
        </Button>
      </div>
      <div className="mt-6 flex flex-col items-center gap-4 sm:flex-row sm:items-start">
        <DealCard />
        <KpiCard />
        <SprintCard />
      </div>
    </section>
  )
}
