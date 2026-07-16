import { Link } from 'react-router-dom'
import { Logo } from '@/components/layout/Logo'

interface AuthTopBarProps {
  navLabel: string
  navHref: string
  /** Shorter fallback shown below `sm` — the full label wraps awkwardly next to the logo at narrow widths. */
  navLabelCompact?: string
}

// Full-width bordered bar with a logo and a cross-link to the other auth
// screen ("Sign in" ↔ "Sign up") — matches the Auth section of the mockups,
// which is a distinct treatment from the plain centered wordmark used by the
// onboarding steps.
export function AuthTopBar({ navLabel, navHref, navLabelCompact }: AuthTopBarProps) {
  return (
    <div className="flex items-center justify-between gap-3 border-b border-border bg-white px-6 py-4 sm:px-20">
      <Logo />
      <Link
        to={navHref}
        className="whitespace-nowrap text-sm text-text-secondary hover:text-foreground"
      >
        {navLabelCompact ? (
          <>
            <span className="hidden sm:inline">{navLabel}</span>
            <span className="sm:hidden">{navLabelCompact}</span>
          </>
        ) : (
          navLabel
        )}
      </Link>
    </div>
  )
}
