import { Link } from 'react-router-dom'

export function Footer() {
  return (
    <footer className="flex items-center justify-between bg-[#050505] px-6 py-5 sm:px-10 lg:px-[100px]">
      <span className="text-xs text-[#444444]">© 2026 FoundryHQ</span>
      <div className="flex items-center gap-6">
        <Link to="/privacy" className="text-xs text-[#444444] hover:text-white">
          Privacy
        </Link>
        <Link to="/terms" className="text-xs text-[#444444] hover:text-white">
          Terms
        </Link>
        <a href="mailto:hello@foundryhq.com" className="text-xs text-[#444444] hover:text-white">
          Contact
        </a>
      </div>
    </footer>
  )
}
