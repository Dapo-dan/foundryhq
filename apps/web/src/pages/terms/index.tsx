import { LegalPageShell } from '@/components/legal/LegalPageShell'

export function TermsPage() {
  return (
    <LegalPageShell title="Terms of Service">
      <p>
        We're still writing this page. FoundryHQ's full terms of service will be published here
        before general availability.
      </p>
      <p>
        In the meantime, if you have questions, reach out at{' '}
        <a href="mailto:hello@foundryhq.com" className="underline hover:text-foreground">
          hello@foundryhq.com
        </a>
        .
      </p>
    </LegalPageShell>
  )
}
