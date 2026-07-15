import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { RadioGroup } from '@/components/ui/radio-group'
import { OptionCard } from '@/pages/onboarding/components/OptionCard'
import { useOnboardingStore } from '@/store/slices/onboarding'
import type { Role } from '@/types/onboarding'

const OPTIONS: { value: Role; title: string; description: string }[] = [
  { value: 'founder_ceo', title: 'Founder / CEO', description: 'Running the whole show' },
  { value: 'coo_operator', title: 'COO / Operator', description: 'Keeping things running smoothly' },
  { value: 'head_of_sales', title: 'Head of Sales', description: 'Driving revenue and growth' },
  { value: 'product_engineer', title: 'Product / Engineer', description: 'Building the product' },
]

export function RoleStepPage() {
  const navigate = useNavigate()
  const storedRole = useOnboardingStore((state) => state.role)
  const setRole = useOnboardingStore((state) => state.setRole)
  const markStepComplete = useOnboardingStore((state) => state.markStepComplete)
  const [selected, setSelected] = useState<Role | null>(storedRole)

  function onContinue() {
    if (!selected) return
    setRole(selected)
    markStepComplete('role')
    navigate('/onboarding/tools')
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1 text-center">
        <h1 className="text-2xl font-bold text-foreground">What's your primary role?</h1>
        <p className="text-sm text-muted-foreground">
          We'll personalize your workspace around what you do.
        </p>
      </div>
      <RadioGroup
        value={selected ?? undefined}
        onValueChange={(value) => setSelected(value as Role)}
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
