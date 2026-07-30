import type { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { useNavigation } from '@react-navigation/native';
import { signUpSchema } from '@foundryhq/shared-validation';
import { useState } from 'react';
import { ScrollView, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { AuthCard } from '../../components/AuthCard';
import { AuthHeader } from '../../components/AuthHeader';
import { PasswordField } from '../../components/PasswordField';
import { PrimaryButton } from '../../components/PrimaryButton';
import { TextField } from '../../components/TextField';
import { useSignUp } from '../../hooks/useSignUp';
import type { AuthStackParamList } from '../../navigation/types';

type Navigation = NativeStackNavigationProp<AuthStackParamList, 'SignUp'>;

export function SignUpScreen() {
  const navigation = useNavigation<Navigation>();
  const signUp = useSignUp();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({});

  function onSubmit() {
    const result = signUpSchema.safeParse({ email, password });
    if (!result.success) {
      const fieldErrors = result.error.flatten().fieldErrors;
      setErrors({ email: fieldErrors.email?.[0], password: fieldErrors.password?.[0] });
      return;
    }
    setErrors({});
    // On success, useSignUp's onSuccess writes the session to the auth
    // store; RootNavigator sees accessToken set + onboardingComplete false
    // and mounts OnboardingNavigator automatically.
    signUp.mutate(result.data);
  }

  return (
    // AuthHeader handles the top safe-area edge itself — only bottom is left
    // to this SafeAreaView, otherwise the two would double up on top padding.
    <SafeAreaView edges={['bottom']} className="flex-1 bg-white">
      <ScrollView className="flex-1">
        <AuthHeader navLabel="Sign in" onNavPress={() => navigation.navigate('SignIn')} />
        <View className="flex-1 items-center justify-center gap-4 px-6 py-12">
          <View className="w-full max-w-[440px]">
            <AuthCard
              heading="Create your account"
              description="14-day free trial. No credit card required."
            >
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
                  autoComplete="new-password"
                  value={password}
                  onChangeText={setPassword}
                  error={errors.password}
                />

                {signUp.isError && (
                  <Text className="text-sm text-red-600">{signUp.error.message}</Text>
                )}

                <PrimaryButton
                  label="Create account"
                  onPress={onSubmit}
                  loading={signUp.isPending}
                />
              </View>
            </AuthCard>
          </View>
          <Text className="max-w-[420px] text-center text-xs text-text-subtle">
            By signing up you agree to our Terms of Service and Privacy Policy
          </Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
