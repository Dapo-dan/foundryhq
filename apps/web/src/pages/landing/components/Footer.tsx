const FOOTER_LINKS = ['Privacy', 'Terms', 'Contact']

export function Footer() {
  return (
    <footer className="flex items-center justify-between bg-[#050505] px-6 py-5 sm:px-10 lg:px-[100px]">
      <span className="text-xs text-[#444444]">© 2026 FoundryHQ</span>
      <div className="flex items-center gap-6">
        {FOOTER_LINKS.map((label) => (
          <a key={label} href="#" className="text-xs text-[#444444] hover:text-white">
            {label}
          </a>
        ))}
      </div>
    </footer>
  )
}
