import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Reveal } from './Reveal'
import { SectionContainer } from './SectionContainer'

export function FinalCta() {
  return (
    <SectionContainer className="flex flex-col items-center bg-lp-dark py-20 text-center">
      <Reveal className="flex flex-col items-center gap-6">
        <h2 className="text-4xl font-extrabold text-white sm:text-5xl lg:text-[52px]">
          From day zero to Series B.
        </h2>
        <p className="max-w-lg text-[17px] leading-relaxed text-[#666666]">
          One workspace that grows with your startup. No bloat. No context-switching. Just
          momentum.
        </p>
        <Button asChild size="lg" className="bg-white px-10 text-base text-lp-dark hover:bg-white/90">
          <Link to="/onboarding">Start building today →</Link>
        </Button>
        <p className="text-[13px] text-[#444444]">Free to start · No credit card required</p>
      </Reveal>
    </SectionContainer>
  )
}
