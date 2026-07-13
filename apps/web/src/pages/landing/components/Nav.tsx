import { Button } from '@/components/ui/button'

const NAV_LINKS = ['Features', 'CRM', 'Roadmap', 'Pricing']

export function Nav() {
  return (
    <nav className="flex items-center justify-between border-b border-border px-6 py-5 sm:px-10 lg:px-[100px]">
      <div className="flex items-center gap-2">
        <div className="h-[26px] w-[26px] rounded bg-brand" />
        <span className="text-[15px] font-bold text-text-primary">FoundryHQ</span>
      </div>
      <div className="hidden items-center gap-9 md:flex">
        {NAV_LINKS.map((label) => (
          <a key={label} href="#" className="text-sm text-text-secondary hover:text-text-primary">
            {label}
          </a>
        ))}
      </div>
      <div className="flex items-center gap-5">
        <a href="#" className="hidden text-sm text-text-secondary hover:text-text-primary sm:inline">
          Sign in
        </a>
        <Button>Get started</Button>
      </div>
    </nav>
  )
}
