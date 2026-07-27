// No red/amber/green tokens exist in this design system (only `destructive`
// for errors) — strength is conveyed by fill amount + label instead of by
// traffic-light color, to stay within the semantic token set.
export function getPasswordStrengthScore(password: string) {
  if (!password) return 0

  let score = 0
  if (password.length >= 8) score++
  if (password.length >= 12) score++
  if (/[0-9]/.test(password) && /[a-zA-Z]/.test(password)) score++
  if (/[^a-zA-Z0-9]/.test(password)) score++

  return Math.min(score, 4)
}
