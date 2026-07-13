const LOGOS = ['Buildwise', 'NovaSeed', 'StackForge', 'Luma Labs', 'Quartex', 'Meriton']

export function SocialProof() {
  return (
    <section className="flex flex-col items-center gap-5 bg-surface-bg px-6 py-12 sm:px-10 lg:px-[100px]">
      <p className="text-sm tracking-wide text-text-subtle">Trusted by 500+ early-stage startups</p>
      <div className="flex flex-wrap items-center justify-center gap-x-10 gap-y-4 sm:gap-x-[52px]">
        {LOGOS.map((name) => (
          <span key={name} className="text-sm font-bold text-[#CCCCCC]">
            {name}
          </span>
        ))}
      </div>
    </section>
  )
}
