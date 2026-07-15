import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { useOnboardingStore } from '@/store/slices/onboarding'

const TOOLS = ['Notion', 'Trello', 'Asana', 'Airtable', 'Jira', 'Spreadsheets']

export function ToolsStepPage() {
  const navigate = useNavigate()
  const tools = useOnboardingStore((state) => state.tools)
  const toggleTool = useOnboardingStore((state) => state.toggleTool)
  const markStepComplete = useOnboardingStore((state) => state.markStepComplete)

  function onContinue() {
    markStepComplete('tools')
    navigate('/onboarding/invite')
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1 text-center">
        <h1 className="text-2xl font-bold text-foreground">What are you replacing?</h1>
        <p className="text-sm text-muted-foreground">
          Select all that apply — we'll help you migrate.
        </p>
      </div>
      <div className="grid grid-cols-3 gap-3">
        {TOOLS.map((tool) => {
          const selected = tools.includes(tool)
          return (
            <button
              key={tool}
              type="button"
              onClick={() => toggleTool(tool)}
              className={cn(
                'rounded-lg border px-4 py-3 text-sm font-medium transition-colors',
                selected
                  ? 'border-primary bg-primary text-primary-foreground'
                  : 'border-input text-foreground hover:bg-muted/40'
              )}
            >
              {tool}
            </button>
          )
        })}
      </div>
      <Button type="button" className="h-11 w-full text-[15px]" onClick={onContinue}>
        Continue →
      </Button>
    </div>
  )
}
