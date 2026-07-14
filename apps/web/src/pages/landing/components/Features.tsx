import { BarChart3, Users, Zap } from 'lucide-react'
import { SectionContainer } from './SectionContainer'

const FEATURES = [
  {
    icon: Users,
    title: 'Lightweight CRM',
    desc: 'Track contacts, companies, and deals in a clean pipeline. Activity timelines, notes, and deal stages — without the Salesforce complexity.',
    bg: 'bg-[#EFF6FF]',
  },
  {
    icon: Zap,
    title: 'Sprint & Task Management',
    desc: 'Plan sprints, assign tasks, and track velocity. Connect product work directly to company goals so nothing gets lost in translation.',
    bg: 'bg-[#F0FDF4]',
  },
  {
    icon: BarChart3,
    title: 'KPI Dashboard',
    desc: 'See your whole company in real time — pipeline, MRR, sprint progress, and OKRs. One view, always current, shareable with investors.',
    bg: 'bg-[#FFF7ED]',
  },
]

export function Features() {
  return (
    <SectionContainer
      id="features"
      className="flex scroll-mt-24 flex-col items-center gap-14 py-16 lg:py-24"
    >
      <div className="flex flex-col items-center gap-3.5 text-center">
        <h2 className="max-w-2xl text-3xl font-extrabold text-text-primary sm:text-4xl lg:text-[44px]">
          Everything your team needs
        </h2>
        <p className="max-w-md text-[17px] leading-relaxed text-text-secondary">
          Replace 10 tools with one workspace that actually knows what a startup needs.
        </p>
      </div>
      <div className="grid w-full grid-cols-1 gap-6 md:grid-cols-3">
        {FEATURES.map(({ icon: Icon, title, desc, bg }) => (
          <div key={title} className={`flex flex-col gap-4 rounded-2xl ${bg} p-9`}>
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white">
              <Icon aria-hidden="true" className="h-6 w-6 text-text-primary" />
            </div>
            <h3 className="text-xl font-bold text-text-primary">{title}</h3>
            <p className="text-sm leading-[1.7] text-text-secondary">{desc}</p>
          </div>
        ))}
      </div>
    </SectionContainer>
  )
}
