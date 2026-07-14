import { SectionContainer } from './SectionContainer'

const LOGOS = ['Buildwise', 'NovaSeed', 'StackForge', 'Luma Labs', 'Quartex', 'Meriton']

export function SocialProof() {
  return (
    <SectionContainer className="flex flex-col items-center gap-5 bg-surface-bg py-12">
      <p className="text-sm tracking-wide text-text-subtle">Trusted by 500+ early-stage startups</p>
      <div className="flex flex-wrap items-center justify-center gap-x-10 gap-y-4 sm:gap-x-[52px]">
        {LOGOS.map((name) => (
          <span
            key={name}
            className="text-sm font-bold text-[#CCCCCC] transition-colors hover:text-text-muted"
          >
            {name}
          </span>
        ))}
      </div>
    </SectionContainer>
  )
}
