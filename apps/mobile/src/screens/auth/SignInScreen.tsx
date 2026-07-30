import type { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { useNavigation } from '@react-navigation/native';
import { signInSchema } from '@foundryhq/shared-validation';
import { useState } from 'react';
import { Pressable, ScrollView, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { AuthCard } from '../../components/AuthCard';
import { AuthHeader } from '../../components/AuthHeader';
import { PasswordField } from '../../components/PasswordField';
import { PrimaryButton } from '../../components/PrimaryButton';
import { TextField } from '../../components/TextField';
import { useSignIn } from '../../hooks/useSignIn';
import type { AuthStackParamList } from '../../navigation/types';

type Navigation = NativeStackNavigationProp<AuthStackParamList, 'SignIn'>;

export function SignInScreen() {
  const navigation = useNavigation<Navigation>();
  const signIn = useSignIn();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({});

  function onSubmit() {
    const result = signInSchema.safeParse({ email, password });
    if (!result.success) {
      const fieldErrors = result.error.flatten().fieldErrors;
      setErrors({ email: fieldErrors.email?.[0], password: fieldErrors.password?.[0] });
      return;
    }
    setErrors({});
    // On success, useSignIn's onSuccess writes the session to the auth
    // store, and RootNavigator's state-driven switch takes it from there —
    // no explicit navigate() needed for the top-level transition.
    signIn.mutate(result.data);
  }

  return (
    // AuthHeader handles the top safe-area edge itself — only bottom is left
    // to this SafeAreaView, otherwise the two would double up on top padding.
    <SafeAreaView edges={['bottom']} className="flex-1 bg-white">
      <ScrollView className="flex-1">
        <AuthHeader navLabel="Sign up" onNavPress={() => navigation.navigate('SignUp')} />
        <View className="flex-1 items-center justify-center px-6 py-12">
          <View className="w-full max-w-[440px]">
            <AuthCard heading="Welcome back" description="Sign in to your FoundryHQ workspace.">
              <View className="gap-3">
                <TextField
                  label="Work email"
                  placeholder="you@company.com"
                  autoComplete="email"
                  autoCapitalize="none"
                  keyboardType="email-address"
                  value={email}
                  onChangeText={setEmail}
                  error={errors.email}
                />
                <PasswordField
                  label="Password"
                  placeholder="At least 8 characters"
                  autoComplete="current-password"
                  value={password}
                  onChangeText={setPassword}
                  error={errors.password}
                />
                <Pressable onPress={() => navigation.navigate('ForgotPassword')}>
                  <Text className="text-xs text-brand">Forgot password?</Text>
                </Pressable>

                {signIn.isError && (
                  <Text className="text-sm text-red-600">{signIn.error.message}</Text>
                )}

                <PrimaryButton label="Sign in" onPress={onSubmit} loading={signIn.isPending} />
              </View>
            </AuthCard>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
