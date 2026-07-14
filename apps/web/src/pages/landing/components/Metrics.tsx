import { useCountUp } from '@/hooks/useCountUp'
import { useInView } from '@/hooks/useInView'
import { SectionContainer } from './SectionContainer'

const METRICS = [
  { target: 10, suffix: '+', label: 'Tools replaced' },
  { target: 2, suffix: 'h', label: 'Saved per week, per person' },
  { target: 500, suffix: '+', label: 'Early-stage startups' },
  { target: 0, prefix: 'Day ', label: 'Time to first value' },
]

function MetricStat({
  target,
  prefix = '',
  suffix = '',
  label,
  active,
}: {
  target: number
  prefix?: string
  suffix?: string
  label: string
  active: boolean
}) {
  const value = useCountUp(target, active)
  return (
    <div className="flex flex-col items-center gap-1.5 text-center">
      <span className="text-4xl font-extrabold text-white sm:text-5xl lg:text-[52px]">
        {prefix}
        {value}
        {suffix}
      </span>
      <span className="max-w-[160px] text-sm text-[#94A3B8]">{label}</span>
    </div>
  )
}

export function Metrics() {
  const { ref, inView } = useInView<HTMLDivElement>(0.3)

  return (
    <SectionContainer className="flex justify-center bg-[#0F172A] py-16">
      <div
        ref={ref}
        className="flex w-full max-w-[1240px] flex-col items-center justify-center gap-10 sm:flex-row sm:justify-between"
      >
        {METRICS.map((metric) => (
          <MetricStat key={metric.label} {...metric} active={inView} />
        ))}
      </div>
    </SectionContainer>
  )
}
