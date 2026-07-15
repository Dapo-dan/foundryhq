import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { inviteEmailSchema } from '@/lib/validation/onboarding'
import { useOnboardingStore } from '@/store/slices/onboarding'
import { InviteEmailField } from './components/InviteEmailField'

export function InviteStepPage() {
  const navigate = useNavigate()
  const storedInvites = useOnboardingStore((state) => state.invites)
  const setInvites = useOnboardingStore((state) => state.setInvites)
  const markStepComplete = useOnboardingStore((state) => state.markStepComplete)

  const [emails, setEmails] = useState<string[]>(
    storedInvites.length ? storedInvites : ['', '', '']
  )
  const [errors, setErrors] = useState<Record<number, string>>({})

  function updateEmail(index: number, value: string) {
    setEmails((prev) => prev.map((email, i) => (i === index ? value : email)))
  }

  function addAnother() {
    setEmails((prev) => [...prev, ''])
  }

  function validate(): boolean {
    const nextErrors: Record<number, string> = {}
    emails.forEach((email, i) => {
      if (email && !inviteEmailSchema.safeParse(email).success) {
        nextErrors[i] = 'Enter a valid email address'
      }
    })
    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  function finishStep() {
    markStepComplete('invite')
    navigate('/onboarding/welcome')
  }

  function onSendInvite() {
    if (!validate()) return
    setInvites(emails.filter(Boolean))
    finishStep()
  }

  function onSkip() {
    finishStep()
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1 text-center">
        <h1 className="text-2xl font-bold text-foreground">Invite your team</h1>
        <p className="text-sm text-muted-foreground">
          Work is better together — invite a few teammates to get started.
        </p>
      </div>
      <div className="flex flex-col gap-3">
        {emails.map((email, i) => (
          <InviteEmailField
            key={i}
            value={email}
            error={errors[i]}
            onChange={(value) => updateEmail(i, value)}
          />
        ))}
        <button
          type="button"
          onClick={addAnother}
          className="self-start text-sm text-brand hover:text-brand-accent"
        >
          + Add another
        </button>
      </div>
      <Button type="button" className="h-11 w-full text-[15px]" onClick={onSendInvite}>
        Send invite →
      </Button>
      <button
        type="button"
        onClick={onSkip}
        className="text-center text-sm text-muted-foreground hover:text-foreground"
      >
        Skip for now
      </button>
    </div>
  )
}
