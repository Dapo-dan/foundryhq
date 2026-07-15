import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { SectionContainer } from './SectionContainer'

const FOCUS_RING = 'rounded-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-accent'

export function Nav() {
  return (
    <SectionContainer
      as="nav"
      className="sticky top-0 z-10 flex items-center justify-between border-b border-border bg-white/80 py-5 backdrop-blur"
    >
      <div className="flex items-center gap-2">
        <div className="h-[26px] w-[26px] rounded bg-brand" />
        <span className="text-[15px] font-bold text-text-primary">FoundryHQ</span>
      </div>
      <div className="flex items-center gap-9">
        <a href="#features" className={`text-sm text-text-secondary hover:text-text-primary ${FOCUS_RING}`}>
          Features
        </a>
      </div>
      <div className="flex items-center gap-5">
        <Link
          to="/auth/sign-in"
          className={`hidden text-sm text-text-secondary hover:text-text-primary sm:inline ${FOCUS_RING}`}
        >
          Sign in
        </Link>
        <Button asChild>
          <Link to="/onboarding">Get started</Link>
        </Button>
      </div>
    </SectionContainer>
  )
}
