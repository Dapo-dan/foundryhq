import { LegalPageShell } from '@/components/legal/LegalPageShell'

export function PrivacyPage() {
  return (
    <LegalPageShell title="Privacy Policy">
      <p>
        We're still writing this page. FoundryHQ's full privacy policy will be published here
        before general availability.
      </p>
      <p>
        In the meantime, if you have questions about how we handle your data, reach out at{' '}
        <a href="mailto:hello@foundryhq.com" className="underline hover:text-foreground">
          hello@foundryhq.com
        </a>
        .
      </p>
    </LegalPageShell>
  )
}
