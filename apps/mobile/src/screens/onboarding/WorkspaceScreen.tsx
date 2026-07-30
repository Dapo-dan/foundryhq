import type { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { useNavigation } from '@react-navigation/native';
import { workspaceSchema } from '@foundryhq/shared-validation';
import { useState } from 'react';
import { ScrollView, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { PrimaryButton } from '../../components/PrimaryButton';
import { StepProgressBar } from '../../components/StepProgressBar';
import { TextField } from '../../components/TextField';
import type { OnboardingStackParamList } from '../../navigation/types';
import { useOnboardingStore } from '../../store/slices/onboarding';

type Navigation = NativeStackNavigationProp<OnboardingStackParamList, 'Workspace'>;

export function WorkspaceScreen() {
  const navigation = useNavigation<Navigation>();
  const workspaceName = useOnboardingStore((state) => state.workspaceName);
  const setWorkspaceName = useOnboardingStore((state) => state.setWorkspaceName);
  const markStepComplete = useOnboardingStore((state) => state.markStepComplete);

  const [name, setName] = useState(workspaceName);
  const [error, setError] = useState<string>();

  function onSubmit() {
    const result = workspaceSchema.safeParse({ name });
    if (!result.success) {
      setError(result.error.flatten().fieldErrors.name?.[0]);
      return;
    }
    setError(undefined);
    // No API call — apps/api has no workspace endpoint yet, matching web's
    // current behavior (local state only, see docs behind this port's plan).
    setWorkspaceName(result.data.name);
    markStepComplete('workspace');
    navigation.navigate('Invite');
  }

  return (
    // No separate header component on this screen (unlike auth screens'
    // AuthHeader) — this SafeAreaView needs both edges.
    <SafeAreaView edges={['top', 'bottom']} className="flex-1 bg-white">
      <ScrollView className="flex-1" contentContainerClassName="flex-1 justify-center px-6 py-12">
        <View className="w-full max-w-[440px] gap-6 self-center">
          <StepProgressBar currentStep={1} totalSteps={2} />
          <View className="gap-1">
            <Text className="text-center text-2xl font-bold text-text-primary">
              Name your workspace
            </Text>
            <Text className="text-center text-sm text-text-secondary">
              You can always change this later.
            </Text>
          </View>
          <View className="gap-3">
            <TextField
              placeholder="e.g. Acme Inc."
              autoFocus
              value={name}
              onChangeText={setName}
              error={error}
            />
            <PrimaryButton label="Continue" onPress={onSubmit} />
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
