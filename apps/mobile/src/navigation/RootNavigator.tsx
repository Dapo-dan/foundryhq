import { ActivityIndicator, View } from 'react-native';
import { useAuthStore } from '../store/slices/auth';
import { useOnboardingStore } from '../store/slices/onboarding';
import { AuthNavigator } from './AuthNavigator';
import { MainTabNavigator } from './MainTabNavigator';
import { OnboardingNavigator } from './OnboardingNavigator';

// Web has no route protection at all (any URL renders regardless of auth
// state) — that works there because there's no path a signed-out user could
// stumble onto by "just not clicking a link." Mobile has no such escape
// hatch, so which tree mounts here is a real, load-bearing decision.
export function RootNavigator() {
  const isHydrated = useAuthStore((state) => state.isHydrated);
  const accessToken = useAuthStore((state) => state.accessToken);
  const onboardingComplete = useOnboardingStore((state) => state.onboardingComplete);

  if (!isHydrated) {
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
