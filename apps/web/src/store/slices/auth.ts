import { create } from 'zustand'
import type { AuthSession, User } from '@foundryhq/shared-types'

interface AuthState {
  user: User | null
  accessToken: string | null
  // True until the initial session-bootstrap request (see
  // hooks/useSessionBootstrap.ts) resolves. AppLayout waits on this before
  // redirecting to sign-in, so a page refresh doesn't bounce a signed-in
  // user for the instant before the refresh-token cookie restores them.
  isBootstrapping: boolean
  setSession: (session: AuthSession) => void
  clearSession: () => void
  finishBootstrapping: () => void
}

// Deliberately not persisted (no `persist` middleware): per ADR-0004 the
// access token must live in memory only on web, never localStorage/sessionStorage.
// A page refresh is expected to require a silent refresh via the httpOnly
// refresh-token cookie, not a read from this store.
export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  accessToken: null,
  isBootstrapping: true,
  setSession: (session) => set({ user: session.user, accessToken: session.accessToken }),
  clearSession: () => set({ user: null, accessToken: null }),
  finishBootstrapping: () => set({ isBootstrapping: false }),
}))
