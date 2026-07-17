import axios, { AxiosError } from 'axios'
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

apiClient.interceptors.response.use(
  // Cast needed: axios types this handler as returning `AxiosResponse`, but
  // we deliberately unwrap to the inner payload here (see note below).
  (response) => (response.data as ApiEnvelope<unknown>).data as unknown as typeof response,
  (error: AxiosError<ApiErrorEnvelope>) => {
    const apiError: ApiError = error.response?.data?.error ?? {
      code: 'network_error',
      message: error.message,
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
