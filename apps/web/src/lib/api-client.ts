import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios'
import type { AuthSession } from '@foundryhq/shared-types'
import { useAuthStore } from '@/store/slices/auth'

export interface ApiError {
  code: string
  message: string
  field?: string
}

interface ApiEnvelope<T> {
  data: T
}

interface ApiErrorEnvelope {
  error: ApiError
}

// Marks a request that's already gone through one 401-retry attempt, so the
// interceptor below never retries the same request twice.
type RetriableRequestConfig = InternalAxiosRequestConfig & { _retriedAfterRefresh?: boolean }

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080',
  // Refresh token is an httpOnly cookie (ADR-0004) — the browser needs to
  // send/receive it on cross-origin requests in dev (different Vite/API ports).
  withCredentials: true,
})

apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Shared across concurrent 401s so a burst of requests triggers one
// /auth/refresh call, not one per request.
let refreshPromise: Promise<AuthSession> | null = null

function refreshAccessToken(): Promise<AuthSession> {
  refreshPromise ??= apiClient
    .post('/auth/refresh')
    .then((response) => response as unknown as AuthSession)
    .finally(() => {
      refreshPromise = null
    })
  return refreshPromise
}

apiClient.interceptors.response.use(
  // Cast needed: axios types this handler as returning `AxiosResponse`, but
  // we deliberately unwrap to the inner payload here (see note below).
  (response) => (response.data as ApiEnvelope<unknown>).data as unknown as typeof response,
  async (error: AxiosError<ApiErrorEnvelope>) => {
    const originalRequest = error.config as RetriableRequestConfig | undefined
    const isAuthEndpoint = originalRequest?.url?.startsWith('/auth/')
    const apiError: ApiError = error.response?.data?.error ?? {
      code: 'network_error',
      message: error.message,
    }

    // On an expired access token, attempt one silent refresh (via the
    // httpOnly refresh-token cookie) and retry the original request exactly
    // once. /auth/* requests are excluded — a 401 there is a real
    // credential/session failure, not an expired access token to recover
    // from.
    if (
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retriedAfterRefresh &&
      !isAuthEndpoint
    ) {
      originalRequest._retriedAfterRefresh = true
      try {
        const session = await refreshAccessToken()
        useAuthStore.getState().setSession(session)
        return apiClient(originalRequest)
      } catch {
        useAuthStore.getState().clearSession()
        window.location.href = '/auth/sign-in'
      }
    }

    return Promise.reject(apiError)
  }
)

// The response interceptor above unwraps `{ data }` at runtime, so every
// request actually resolves to `T`, not `AxiosResponse<T>` — axios's own
// types can't express that, hence the cast. Services should call these
// instead of `apiClient.get/post/...` directly.
export function apiGet<TResponse>(url: string) {
  return apiClient.get(url) as unknown as Promise<TResponse>
}

export function apiPost<TResponse, TBody = unknown>(url: string, body?: TBody) {
  return apiClient.post(url, body) as unknown as Promise<TResponse>
}

export function apiPatch<TResponse, TBody = unknown>(url: string, body?: TBody) {
  return apiClient.patch(url, body) as unknown as Promise<TResponse>
}

export function apiPut<TResponse, TBody = unknown>(url: string, body?: TBody) {
  return apiClient.put(url, body) as unknown as Promise<TResponse>
}

export function apiDelete<TResponse>(url: string) {
  return apiClient.delete(url) as unknown as Promise<TResponse>
}
