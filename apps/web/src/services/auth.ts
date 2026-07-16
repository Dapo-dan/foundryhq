import { apiPost } from '@/lib/api-client'
import { USE_MOCK_API } from '@/lib/mock/config'
import { mockForgotPassword, mockResetPassword, mockSignIn, mockSignUp } from '@/lib/mock/auth'
import type {
  AuthSession,
  ForgotPasswordInput,
  ResetPasswordInput,
  SignInInput,
  SignUpInput,
} from '@/types/auth'

export function signUp(input: SignUpInput) {
  return USE_MOCK_API ? mockSignUp(input) : apiPost<AuthSession>('/auth/register', input)
}

export function signIn(input: SignInInput) {
  return USE_MOCK_API ? mockSignIn(input) : apiPost<AuthSession>('/auth/login', input)
}

// NOTE: not in docs/api.md yet — the endpoint index only lists
// register/login/oauth/refresh/logout. Shape follows the same envelope/verb
// conventions as the rest of the Auth section; confirm against the real
// handler once it exists.
export function forgotPassword(input: ForgotPasswordInput) {
  return USE_MOCK_API
    ? mockForgotPassword(input)
    : apiPost<void>('/auth/forgot-password', input)
}

// NOTE: also undocumented — see forgotPassword above.
export function resetPassword(input: ResetPasswordInput) {
  return USE_MOCK_API ? mockResetPassword(input) : apiPost<void>('/auth/reset-password', input)
}
