export interface User {
  id: string
  email: string
}

export interface AuthSession {
  user: User
  accessToken: string
}

export interface SignUpInput {
  email: string
  password: string
  marketingOptIn: boolean
}

export interface SignInInput {
  email: string
  password: string
}

export interface ForgotPasswordInput {
  email: string
}

export interface ResetPasswordInput {
  token: string
  password: string
}
