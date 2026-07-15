import { cn } from '@/lib/utils'

interface PasswordStrengthBarProps {
  password: string
}

const LABELS = ['Weak', 'Fair', 'Good', 'Strong']

// No red/amber/green tokens exist in this design system (only `destructive`
// for errors) — strength is conveyed by fill amount + label instead of by
// traffic-light color, to stay within the semantic token set.
function getScore(password: string) {
  if (!password) return 0

  let score = 0
  if (password.length >= 8) score++
  if (password.length >= 12) score++
  if (/[0-9]/.test(password) && /[a-zA-Z]/.test(password)) score++
  if (/[^a-zA-Z0-9]/.test(password)) score++

  return Math.min(score, 4)
}

export function PasswordStrengthBar({ password }: PasswordStrengthBarProps) {
  const score = getScore(password)

  if (!password) return null

  return (
    <div className="flex items-center gap-2">
      <div className="flex flex-1 gap-1">
        {Array.from({ length: 4 }, (_, i) => (
          <div
            key={i}
            className={cn('h-1 flex-1 rounded-full', i < score ? 'bg-brand-accent' : 'bg-muted')}
          />
        ))}
      </div>
      <span className="text-xs text-text-subtle">{LABELS[Math.max(score - 1, 0)]}</span>
    </div>
  )
}
