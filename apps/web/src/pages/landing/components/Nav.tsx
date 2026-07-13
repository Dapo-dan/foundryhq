import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'

export function Nav() {
  return (
    <nav className="sticky top-0 z-10 flex items-center justify-between border-b border-border bg-white/80 px-6 py-5 backdrop-blur sm:px-10 lg:px-[100px]">
      <div className="flex items-center gap-2">
        <div className="h-[26px] w-[26px] rounded bg-brand" />
        <span className="text-[15px] font-bold text-text-primary">FoundryHQ</span>
      </div>
      <div className="flex items-center gap-9">
        <a href="#features" className="text-sm text-text-secondary hover:text-text-primary">
          Features
        </a>
      </div>
      <div className="flex items-center gap-5">
        <Link
          to="/auth"
          className="hidden text-sm text-text-secondary hover:text-text-primary sm:inline"
        >
          Sign in
        </Link>
        <Button asChild>
          <Link to="/onboarding">Get started</Link>
        </Button>
      </div>
    </nav>
  )
}
