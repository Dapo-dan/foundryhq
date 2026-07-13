import { Button } from '@/components/ui/button'

export function FinalCta() {
  return (
    <section className="flex flex-col items-center gap-6 bg-lp-dark px-6 py-20 text-center sm:px-10 lg:px-[100px]">
      <h2 className="text-4xl font-extrabold text-white sm:text-5xl lg:text-[52px]">
        From day zero to Series B.
      </h2>
      <p className="max-w-lg text-[17px] leading-relaxed text-[#666666]">
        One workspace that grows with your startup. No bloat. No context-switching. Just momentum.
      </p>
      <Button size="lg" className="bg-white px-10 text-base text-lp-dark hover:bg-white/90">
        Start building today →
      </Button>
      <p className="text-[13px] text-[#444444]">Free to start · No credit card required</p>
    </section>
  )
}
