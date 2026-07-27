import './global.css';

import { NavigationContainer, type LinkingOptions } from '@react-navigation/native';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEffect } from 'react';
import { StatusBar } from 'expo-status-bar';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { RootNavigator } from './src/navigation/RootNavigator';
import type { AuthStackParamList } from './src/navigation/types';
import { useAuthStore } from './src/store/slices/auth';

const queryClient = new QueryClient();

const linking: LinkingOptions<AuthStackParamList> = {
  prefixes: ['foundryhq://'],
  config: {
    screens: {
      SignIn: 'sign-in',
      SignUp: 'sign-up',
      ForgotPassword: 'forgot-password',
      ResetPassword: 'reset-password',
    },
  },
};

export default function App() {
  const hydrate = useAuthStore((state) => state.hydrate);

  useEffect(() => {
    hydrate();
  }, [hydrate]);

  return (
    <QueryClientProvider client={queryClient}>
      <SafeAreaProvider>
        <NavigationContainer linking={linking}>
          <RootNavigator />
        </NavigationContainer>
        <StatusBar style="auto" />
      </SafeAreaProvider>
    </QueryClientProvider>
  );
}
