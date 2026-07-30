import { apiPost } from '@/lib/api-client'
import type {
  AuthSession,
  ForgotPasswordInput,
  ResetPasswordInput,
  SignInInput,
  SignUpInput,
} from '@foundryhq/shared-types'

export function signUp(input: SignUpInput) {
  return apiPost<AuthSession>('/auth/register', input)
}

export function signIn(input: SignInInput) {
  return apiPost<AuthSession>('/auth/login', input)
}

export function forgotPassword(input: ForgotPasswordInput) {
  return apiPost<void>('/auth/forgot-password', input)
}

export function resetPassword(input: ResetPasswordInput) {
  return apiPost<void>('/auth/reset-password', input)
}

// Exchanges the httpOnly refresh-token cookie (ADR-0004) for a fresh
// access token. Used both by useSessionBootstrap (on app load) and by
// api-client.ts's 401 retry interceptor (mid-session).
export function refreshSession() {
  return apiPost<AuthSession>('/auth/refresh')
}

export function logout() {
  return apiPost<void>('/auth/logout')
}
