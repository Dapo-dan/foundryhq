import { useEffect, useRef } from 'react'
import { refreshSession } from '@/services/auth'
import { useAuthStore } from '@/store/slices/auth'

// Runs once on app load, exchanging the httpOnly refresh-token cookie
// (ADR-0004) for a fresh access token so a page refresh doesn't drop the
// session. Doesn't block rendering — public pages and the auth screens
// render immediately regardless of outcome; only AppLayout's guard waits on
// `isBootstrapping` before deciding whether to redirect to sign-in.
export function useSessionBootstrap() {
  const ranRef = useRef(false)

  useEffect(() => {
    if (ranRef.current) return
    ranRef.current = true

    refreshSession()
      .then((session) => useAuthStore.getState().setSession(session))
      .catch(() => {
        // No valid refresh-token cookie — the visitor simply isn't signed in.
      })
      .finally(() => useAuthStore.getState().finishBootstrapping())
  }, [])
}
