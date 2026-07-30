import type { NativeStackNavigationProp, NativeStackScreenProps } from '@react-navigation/native-stack';
import { useNavigation } from '@react-navigation/native';
import { resetPasswordSchema } from '@foundryhq/shared-validation';
import { useState } from 'react';
import { ScrollView, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { AuthCard } from '../../components/AuthCard';
import { AuthHeader } from '../../components/AuthHeader';
import { PasswordField } from '../../components/PasswordField';
import { PasswordStrengthBar } from '../../components/PasswordStrengthBar';
import { PrimaryButton } from '../../components/PrimaryButton';
import { useResetPassword } from '../../hooks/useResetPassword';
import type { AuthStackParamList } from '../../navigation/types';

type Navigation = NativeStackNavigationProp<AuthStackParamList, 'ResetPassword'>;
type Props = NativeStackScreenProps<AuthStackParamList, 'ResetPassword'>;

export function ResetPasswordScreen({ route }: Props) {
  const navigation = useNavigation<Navigation>();
  const { token } = route.params;
  const resetPassword = useResetPassword();
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState<{ password?: string; confirmPassword?: string }>({});

  function onSubmit() {
    const result = resetPasswordSchema.safeParse({ password, confirmPassword });
    if (!result.success) {
      const fieldErrors = result.error.flatten().fieldErrors;
      setErrors({
        password: fieldErrors.password?.[0],
        confirmPassword: fieldErrors.confirmPassword?.[0],
      });
      return;
    }
    setErrors({});
    resetPassword.mutate(
      { token, password: result.data.password },
      { onSuccess: () => navigation.navigate('SignIn') }
    );
  }

  return (
    // AuthHeader handles the top safe-area edge itself — only bottom is left
    // to this SafeAreaView, otherwise the two would double up on top padding.
    <SafeAreaView edges={['bottom']} className="flex-1 bg-white">
      <ScrollView className="flex-1">
        <AuthHeader navLabel="← Back to sign in" onNavPress={() => navigation.navigate('SignIn')} />
        <View className="flex-1 items-center justify-center px-6 py-12">
          <View className="w-full max-w-[440px]">
            <AuthCard
              heading="Set a new password"
              description="Choose a new password for your account."
            >
              <View className="gap-3">
                <View className="gap-1.5">
                  <PasswordField
                    label="New password"
                    placeholder="At least 8 characters"
                    autoComplete="new-password"
                    value={password}
                    onChangeText={setPassword}
                    error={errors.password}
                  />
                  <PasswordStrengthBar password={password} />
                </View>
                <PasswordField
                  label="Confirm password"
                  placeholder="Re-enter your password"
                  autoComplete="new-password"
                  value={confirmPassword}
                  onChangeText={setConfirmPassword}
                  error={errors.confirmPassword}
                />

                {resetPassword.isError && (
                  <Text className="text-sm text-red-600">{resetPassword.error.message}</Text>
                )}

                <PrimaryButton
                  label="Update password"
                  onPress={onSubmit}
                  loading={resetPassword.isPending}
                />
              </View>
            </AuthCard>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
