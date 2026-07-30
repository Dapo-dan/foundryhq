import { Navigate, Outlet } from 'react-router-dom'
import { Sidebar } from '@/components/layout/Sidebar'
import { TopNav } from '@/components/layout/TopNav'
import { useAuthStore } from '@/store/slices/auth'

export function AppLayout() {
  const user = useAuthStore((state) => state.user)
  const isBootstrapping = useAuthStore((state) => state.isBootstrapping)

  // Wait for the initial session-bootstrap request (see
  // hooks/useSessionBootstrap.ts) before deciding — otherwise a page
  // refresh would flash a redirect to sign-in before the httpOnly
  // refresh-token cookie has had a chance to restore the session.
  if (isBootstrapping) {
    return null
  }

  if (!user) {
    return <Navigate to="/auth/sign-in" replace />
  }

  return (
    <div className="flex h-svh">
      <Sidebar />
      <div className="flex min-w-0 flex-1 flex-col">
        <TopNav />
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
