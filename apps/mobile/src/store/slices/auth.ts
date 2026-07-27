import { create } from 'zustand'
import type { AuthSession, User } from '@foundryhq/shared-types'
import { deleteSecureItem, getSecureItem, setSecureItem } from '../../services/secure-storage'

const USER_KEY = 'auth.user'
const ACCESS_TOKEN_KEY = 'auth.accessToken'

interface AuthState {
  user: User | null
  accessToken: string | null
  isHydrated: boolean
  setSession: (session: AuthSession) => Promise<void>
  clearSession: () => Promise<void>
  hydrate: () => Promise<void>
}

// Unlike web (ADR-0004: access token in memory only, refresh token via an
// httpOnly cookie), mobile has no cookie jar — per docs/architecture.md's
// Auth Flow, the access token persists to SecureStore so it survives an app
// restart. `hydrate()` reads it back on boot; RootNavigator waits on
// `isHydrated` before deciding which navigator tree to mount.
export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  accessToken: null,
  isHydrated: false,
  setSession: async (session) => {
    set({ user: session.user, accessToken: session.accessToken })
    await Promise.all([
      setSecureItem(USER_KEY, JSON.stringify(session.user)),
      setSecureItem(ACCESS_TOKEN_KEY, session.accessToken),
    ])
  },
  clearSession: async () => {
    set({ user: null, accessToken: null })
    await Promise.all([deleteSecureItem(USER_KEY), deleteSecureItem(ACCESS_TOKEN_KEY)])
  },
  hydrate: async () => {
    const [userJson, accessToken] = await Promise.all([
      getSecureItem(USER_KEY),
      getSecureItem(ACCESS_TOKEN_KEY),
    ])
    set({
      user: userJson ? (JSON.parse(userJson) as User) : null,
      accessToken,
      isHydrated: true,
    })
  },
}))
