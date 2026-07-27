import { getPasswordStrengthScore } from '@foundryhq/shared-validation'
import { cn } from '@/lib/utils'

interface PasswordStrengthBarProps {
  password: string
}

const LABELS = ['Weak', 'Fair', 'Good', 'Strong']

export function PasswordStrengthBar({ password }: PasswordStrengthBarProps) {
  const score = getPasswordStrengthScore(password)

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
