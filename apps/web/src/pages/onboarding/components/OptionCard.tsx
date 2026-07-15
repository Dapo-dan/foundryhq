import { Label } from '@/components/ui/label'
import { RadioGroupItem } from '@/components/ui/radio-group'
import { cn } from '@/lib/utils'

interface OptionCardProps {
  value: string
  title: string
  description: string
  selected: boolean
}

// Shared selectable-card presentation for the Team Size / Role steps — each
// option is a full-width bordered row rather than a bare radio dot.
export function OptionCard({ value, title, description, selected }: OptionCardProps) {
  return (
    <Label
      htmlFor={value}
      className={cn(
        'flex cursor-pointer items-center justify-between gap-4 rounded-lg border p-4 font-normal',
        selected ? 'border-primary bg-muted/40' : 'border-input'
      )}
    >
      <div className="flex flex-col gap-0.5">
        <span className="text-sm font-medium text-foreground">{title}</span>
        <span className="text-xs text-muted-foreground">{description}</span>
      </div>
      <RadioGroupItem value={value} id={value} />
    </Label>
  )
}
