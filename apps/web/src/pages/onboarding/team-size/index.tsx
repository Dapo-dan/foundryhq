import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { RadioGroup } from '@/components/ui/radio-group'
import { OptionCard } from '@/pages/onboarding/components/OptionCard'
import { useOnboardingStore } from '@/store/slices/onboarding'
import type { TeamSize } from '@/types/onboarding'

const OPTIONS: { value: TeamSize; title: string; description: string }[] = [
  { value: 'just_me', title: 'Just me', description: "Solo founder, for now" },
  { value: 'small', title: '2-10 people', description: 'Small, fast-moving team' },
  { value: 'medium', title: '11-50 people', description: 'Growing team' },
  { value: 'large', title: '50+ people', description: 'Established organization' },
]

export function TeamSizeStepPage() {
  const navigate = useNavigate()
  const storedTeamSize = useOnboardingStore((state) => state.teamSize)
  const setTeamSize = useOnboardingStore((state) => state.setTeamSize)
  const markStepComplete = useOnboardingStore((state) => state.markStepComplete)
  const [selected, setSelected] = useState<TeamSize | null>(storedTeamSize)

  function onContinue() {
    if (!selected) return
    setTeamSize(selected)
    markStepComplete('team-size')
    navigate('/onboarding/role')
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1 text-center">
        <h1 className="text-2xl font-bold text-foreground">How big is your team?</h1>
        <p className="text-sm text-muted-foreground">
          We'll tailor FoundryHQ to your team's size.
        </p>
      </div>
      <RadioGroup
        value={selected ?? undefined}
        onValueChange={(value) => setSelected(value as TeamSize)}
      >
        {OPTIONS.map((option) => (
          <OptionCard key={option.value} {...option} selected={selected === option.value} />
        ))}
      </RadioGroup>
      <Button
        type="button"
        className="h-11 w-full text-[15px]"
        disabled={!selected}
        onClick={onContinue}
      >
        Continue →
      </Button>
    </div>
  )
}
