import type {
  AuthSession,
  ForgotPasswordInput,
  ResetPasswordInput,
  SignInInput,
  SignUpInput,
} from '@foundryhq/shared-types';
import { apiPost } from './api-client';
import { mockForgotPassword, mockResetPassword, mockSignIn, mockSignUp } from './mock/auth';
import { USE_MOCK_API } from './mock/config';

export function signUp(input: SignUpInput) {
  return USE_MOCK_API ? mockSignUp(input) : apiPost<AuthSession>('/auth/register', input);
}

export function signIn(input: SignInInput) {
  return USE_MOCK_API ? mockSignIn(input) : apiPost<AuthSession>('/auth/login', input);
}

// apps/api has no forgot-password/reset-password route yet — real mode will
// 404 until it exists. Shape follows the same envelope/verb conventions as
// the rest of the Auth section; confirm against the real handler once it exists.
export function forgotPassword(input: ForgotPasswordInput) {
  return USE_MOCK_API
    ? mockForgotPassword(input)
    : apiPost<void>('/auth/forgot-password', input);
}

export function resetPassword(input: ResetPasswordInput) {
  return USE_MOCK_API ? mockResetPassword(input) : apiPost<void>('/auth/reset-password', input);
}
