import { Link } from 'react-router-dom'
import { SectionContainer } from './SectionContainer'

const FOCUS_RING = 'rounded-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-accent'

export function Footer() {
  return (
    <SectionContainer as="footer" className="flex items-center justify-between bg-[#050505] py-5">
      <span className="text-xs text-[#444444]">© 2026 FoundryHQ</span>
      <div className="flex items-center gap-6">
        <Link to="/privacy" className={`text-xs text-[#444444] hover:text-white ${FOCUS_RING}`}>
          Privacy
        </Link>
        <Link to="/terms" className={`text-xs text-[#444444] hover:text-white ${FOCUS_RING}`}>
          Terms
        </Link>
        <a
          href="mailto:hello@foundryhq.com"
          className={`text-xs text-[#444444] hover:text-white ${FOCUS_RING}`}
        >
          Contact
        </a>
      </div>
    </SectionContainer>
  )
}
