import { Link } from 'react-router-dom'

export function Logo() {
  return (
    <Link to="/" className="flex items-center gap-2">
      <div className="size-6 rounded bg-brand" />
      <span className="text-[15px] font-bold text-foreground">FoundryHQ</span>
    </Link>
  )
}
