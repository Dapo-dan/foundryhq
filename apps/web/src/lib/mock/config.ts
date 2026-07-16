// apps/api has no auth routes yet (only health checks) — mock mode lets the
// sign-up/sign-in/password flows run against dummy data instead of a failing
// network call. Opt out once the real endpoints exist with VITE_API_MOCKS=false.
export const USE_MOCK_API = import.meta.env.DEV && import.meta.env.VITE_API_MOCKS !== 'false'
