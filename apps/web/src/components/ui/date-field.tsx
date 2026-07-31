import { useState } from 'react'
import { format } from 'date-fns'
import { CalendarIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

interface DateFieldProps {
  value?: string
  onChange: (value: string) => void
  placeholder?: string
}

// A single-date picker built from the calendar/popover primitives, storing
// its value as a 'yyyy-MM-dd' string — the same date-only shape the API
// expects for sprint start/end dates and task due dates. Shared by
// NewSprintDialog (start/end) and NewTaskDialog (due date) rather than
// wiring the popover+calendar combination three separate times.
export function DateField({ value, onChange, placeholder = 'Pick a date' }: DateFieldProps) {
  const [open, setOpen] = useState(false)
  const selected = value ? new Date(`${value}T00:00:00`) : undefined

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          type="button"
          variant="outline"
          className={cn('w-full justify-start font-normal', !value && 'text-muted-foreground')}
        >
          <CalendarIcon size={16} />
          {selected ? format(selected, 'PPP') : placeholder}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0">
        <Calendar
          mode="single"
          selected={selected}
          onSelect={(date) => {
            if (!date) return
            onChange(format(date, 'yyyy-MM-dd'))
            // A single-date picker has nothing left to do once a day is
            // chosen — closing immediately (rather than waiting for an
            // outside click) matches how the rest of this app's dialogs
            // behave and stops the calendar covering fields below it.
            setOpen(false)
          }}
        />
      </PopoverContent>
    </Popover>
  )
}
