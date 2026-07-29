import { useEffect, useState } from 'react';
import { ActivityIndicator, View } from 'react-native';
import { useAuthStore } from '../store/slices/auth';
import { useOnboardingStore } from '../store/slices/onboarding';
import { AuthNavigator } from './AuthNavigator';
import { MainTabNavigator } from './MainTabNavigator';
import { OnboardingNavigator } from './OnboardingNavigator';

// The onboarding store's `persist` middleware rehydrates from AsyncStorage
// on its own schedule, independent of the auth store's explicit hydrate().
// Without waiting on this too, a returning user with onboarding already
// complete could briefly see OnboardingNavigator mount before rehydration
// catches up and flips `onboardingComplete` — exactly the kind of flash this
// navigator exists to avoid.
function useOnboardingHydrated() {
  const [hydrated, setHydrated] = useState(() => useOnboardingStore.persist.hasHydrated());

  useEffect(() => {
    if (hydrated) return;
    return useOnboardingStore.persist.onFinishHydration(() => setHydrated(true));
  }, [hydrated]);

  return hydrated;
}

// Web has no route protection at all (any URL renders regardless of auth
// state) — that works there because there's no path a signed-out user could
// stumble onto by "just not clicking a link." Mobile has no such escape
// hatch, so which tree mounts here is a real, load-bearing decision.
export function RootNavigator() {
  const isAuthHydrated = useAuthStore((state) => state.isHydrated);
  const isOnboardingHydrated = useOnboardingHydrated();
  const accessToken = useAuthStore((state) => state.accessToken);
  const onboardingComplete = useOnboardingStore((state) => state.onboardingComplete);

  if (!isAuthHydrated || !isOnboardingHydrated) {
    return (
      <View className="flex-1 items-center justify-center bg-white">
        <ActivityIndicator />
      </View>
    );
  }

  if (!accessToken) {
    return <AuthNavigator />;
  }

  if (!onboardingComplete) {
    return <OnboardingNavigator />;
  }

  return <MainTabNavigator />;
}
