import { Input } from '@/components/ui/input'

interface InviteEmailFieldProps {
  value: string
  error?: string
  onChange: (value: string) => void
}

export function InviteEmailField({ value, error, onChange }: InviteEmailFieldProps) {
  return (
    <div>
      <Input
        type="email"
        placeholder="teammate@company.com"
        autoComplete="off"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
      {error && <p className="mt-1 text-sm text-destructive">{error}</p>}
    </div>
  )
}
