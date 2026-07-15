import { Link } from 'react-router-dom'
import { Logo } from '@/components/layout/Logo'

interface AuthTopBarProps {
  navLabel: string
  navHref: string
}

// Full-width bordered bar with a logo and a cross-link to the other auth
// screen ("Sign in" ↔ "Sign up") — matches the Auth section of the mockups,
// which is a distinct treatment from the plain centered wordmark used by the
// onboarding steps.
export function AuthTopBar({ navLabel, navHref }: AuthTopBarProps) {
  return (
    <div className="flex items-center justify-between border-b border-border bg-white px-6 py-4 sm:px-20">
      <Logo />
      <Link to={navHref} className="text-sm text-text-secondary hover:text-foreground">
        {navLabel}
      </Link>
    </div>
  )
}
