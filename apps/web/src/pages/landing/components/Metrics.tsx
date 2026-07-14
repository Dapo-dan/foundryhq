import { SectionContainer } from './SectionContainer'

const METRICS = [
  { value: '10+', label: 'Tools replaced' },
  { value: '2h', label: 'Saved per week, per person' },
  { value: '500+', label: 'Early-stage startups' },
  { value: 'Day 0', label: 'Time to first value' },
]

export function Metrics() {
  return (
    <SectionContainer className="flex flex-col items-center justify-center gap-10 bg-[#0F172A] py-16 sm:flex-row sm:justify-between">
      {METRICS.map(({ value, label }) => (
        <div key={label} className="flex flex-col items-center gap-1.5 text-center">
          <span className="text-4xl font-extrabold text-white sm:text-5xl lg:text-[52px]">{value}</span>
          <span className="max-w-[160px] text-sm text-[#94A3B8]">{label}</span>
        </div>
      ))}
    </SectionContainer>
  )
}
