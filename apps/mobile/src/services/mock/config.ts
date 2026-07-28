// apps/api routes /auth/register and /auth/login for real, but has no
// forgot-password/reset-password endpoints yet — mock mode lets all four
// flows run against dummy data during local dev either way. Opt out once
// the remaining endpoints exist with EXPO_PUBLIC_API_MOCKS=false.
export const USE_MOCK_API = __DEV__ && process.env.EXPO_PUBLIC_API_MOCKS !== 'false';
