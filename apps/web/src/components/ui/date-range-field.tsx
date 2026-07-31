import { useState } from 'react'
import { format } from 'date-fns'
import { CalendarIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

interface DateRangeFieldProps {
  startValue?: string
  endValue?: string
  onChangeStart: (value: string) => void
  onChangeEnd: (value: string) => void
  placeholder?: string
}

// A single popover holding one range-mode calendar, storing each end as a
// 'yyyy-MM-dd' string — same shape DateField uses for a single date. Built
// for NewSprintDialog's start/end pair, which used to be two independent
// DateFields (a picker per field, not one range control).
export function DateRangeField({
  startValue,
  endValue,
  onChangeStart,
  onChangeEnd,
  placeholder = 'Pick a date range',
}: DateRangeFieldProps) {
  const [open, setOpen] = useState(false)
  const from = startValue ? new Date(`${startValue}T00:00:00`) : undefined
  const to = endValue ? new Date(`${endValue}T00:00:00`) : undefined

  let label = placeholder
  if (from && to) {
    label = `${format(from, 'PPP')} – ${format(to, 'PPP')}`
  } else if (from) {
    label = `${format(from, 'PPP')} – …`
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          type="button"
          variant="outline"
          className={cn('w-full justify-start font-normal', !from && 'text-muted-foreground')}
        >
          <CalendarIcon size={16} />
          {label}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0">
        <Calendar
          mode="range"
          selected={{ from, to }}
          onSelect={(range) => {
            onChangeStart(range?.from ? format(range.from, 'yyyy-MM-dd') : '')
            onChangeEnd(range?.to ? format(range.to, 'yyyy-MM-dd') : '')
            // Only close once both ends are picked — closing after the
            // first click would hide the calendar before the end date can
            // be chosen.
            if (range?.from && range?.to) setOpen(false)
          }}
        />
      </PopoverContent>
    </Popover>
  )
}
