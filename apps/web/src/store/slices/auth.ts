import { create } from 'zustand'
import type { AuthSession, User } from '@foundryhq/shared-types'

interface AuthState {
  user: User | null
  accessToken: string | null
  setSession: (session: AuthSession) => void
  clearSession: () => void
}

// Deliberately not persisted (no `persist` middleware): per ADR-0004 the
// access token must live in memory only on web, never localStorage/sessionStorage.
// A page refresh is expected to require a silent refresh via the httpOnly
// refresh-token cookie, not a read from this store.
export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  accessToken: null,
  setSession: (session) => set({ user: session.user, accessToken: session.accessToken }),
  clearSession: () => set({ user: null, accessToken: null }),
}))
