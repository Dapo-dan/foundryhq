import type { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { useNavigation } from '@react-navigation/native';
import { forgotPasswordSchema } from '@foundryhq/shared-validation';
import { useState } from 'react';
import { ScrollView, Text, View } from 'react-native';
import { AuthCard } from '../../components/AuthCard';
import { AuthHeader } from '../../components/AuthHeader';
import { PrimaryButton } from '../../components/PrimaryButton';
import { TextField } from '../../components/TextField';
import { useForgotPassword } from '../../hooks/useForgotPassword';
import type { AuthStackParamList } from '../../navigation/types';

type Navigation = NativeStackNavigationProp<AuthStackParamList, 'ForgotPassword'>;

export function ForgotPasswordScreen() {
  const navigation = useNavigation<Navigation>();
  const forgotPassword = useForgotPassword();
  const [email, setEmail] = useState('');
  const [error, setError] = useState<string>();
  const [sentTo, setSentTo] = useState<string | null>(null);

  function onSubmit() {
    const result = forgotPasswordSchema.safeParse({ email });
    if (!result.success) {
      setError(result.error.flatten().fieldErrors.email?.[0]);
      return;
    }
    setError(undefined);
    forgotPassword.mutate(result.data, { onSuccess: () => setSentTo(result.data.email) });
  }

  return (
    <ScrollView className="flex-1 bg-white">
      <AuthHeader navLabel="← Back to sign in" onNavPress={() => navigation.navigate('SignIn')} />
      <View className="flex-1 items-center justify-center px-6 py-12">
        <View className="w-full max-w-[440px]">
          <AuthCard
            heading="Reset your password"
            description="Enter your email address and we'll send you a reset link."
          >
            {sentTo ? (
              <Text className="text-center text-sm text-text-secondary">
                If an account exists for{' '}
                <Text className="font-medium text-text-primary">{sentTo}</Text>, we've sent a
                password reset link to that address.
              </Text>
            ) : (
              <View className="gap-3">
                <TextField
                  label="Your email"
                  placeholder="you@company.com"
                  autoComplete="email"
                  autoCapitalize="none"
                  keyboardType="email-address"
                  value={email}
                  onChangeText={setEmail}
                  error={error}
                />

                {forgotPassword.isError && (
                  <Text className="text-sm text-red-600">
                    {forgotPassword.error.message}
                  </Text>
                )}

                <PrimaryButton
                  label="Send reset link"
                  onPress={onSubmit}
                  loading={forgotPassword.isPending}
                />
              </View>
            )}
          </AuthCard>
        </View>
      </View>
    </ScrollView>
  );
}
