import { Outlet } from 'react-router-dom'

export function AuthLayout() {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center bg-background px-4">
      <div className="mb-8 text-lg font-semibold">FoundryHQ</div>
      <div className="w-full max-w-sm">
        <Outlet />
      </div>
    </div>
  )
}
