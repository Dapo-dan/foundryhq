import type {
  AuthSession,
  ForgotPasswordInput,
  ResetPasswordInput,
  SignInInput,
  SignUpInput,
} from '@foundryhq/shared-types'

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function mockSession(email: string): AuthSession {
  return {
    user: { id: crypto.randomUUID(), email },
    accessToken: `mock-token-${crypto.randomUUID()}`,
  }
}

export async function mockSignUp(input: SignUpInput): Promise<AuthSession> {
  await delay(500)
  return mockSession(input.email)
}

export async function mockSignIn(input: SignInInput): Promise<AuthSession> {
  await delay(500)
  return mockSession(input.email)
}

export async function mockForgotPassword(_input: ForgotPasswordInput): Promise<void> {
  await delay(500)
}

export async function mockResetPassword(_input: ResetPasswordInput): Promise<void> {
  await delay(500)
}
