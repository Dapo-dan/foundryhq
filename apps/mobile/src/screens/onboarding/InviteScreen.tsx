import { inviteEmailSchema } from '@foundryhq/shared-validation';
import { useState } from 'react';
import { Pressable, ScrollView, Text, View } from 'react-native';
import { PrimaryButton } from '../../components/PrimaryButton';
import { StepProgressBar } from '../../components/StepProgressBar';
import { TextField } from '../../components/TextField';
import { useOnboardingStore } from '../../store/slices/onboarding';

export function InviteScreen() {
  const storedInvites = useOnboardingStore((state) => state.invites);
  const setInvites = useOnboardingStore((state) => state.setInvites);
  const markStepComplete = useOnboardingStore((state) => state.markStepComplete);
  const markOnboardingComplete = useOnboardingStore((state) => state.markOnboardingComplete);

  const [emails, setEmails] = useState<string[]>(
    storedInvites.length ? storedInvites : ['', '', '']
  );
  const [errors, setErrors] = useState<Record<number, string>>({});

  function updateEmail(index: number, value: string) {
    setEmails((prev) => prev.map((email, i) => (i === index ? value : email)));
  }

  function addAnother() {
    setEmails((prev) => [...prev, '']);
  }

  function validate(): boolean {
    const nextErrors: Record<number, string> = {};
    emails.forEach((email, i) => {
      if (email && !inviteEmailSchema.safeParse(email).success) {
        nextErrors[i] = 'Enter a valid email address';
      }
    });
    setErrors(nextErrors);
    return Object.keys(nextErrors).length === 0;
  }

  // With only Workspace and Invite in this port's onboarding scope,
  // finishing this step always means onboarding is done — RootNavigator's
  // state-driven switch takes it from here to MainTabNavigator, no explicit
  // navigate() needed.
  function finishStep() {
    markStepComplete('invite');
    markOnboardingComplete();
  }

  function onSendInvite() {
    if (!validate()) return;
    setInvites(emails.filter(Boolean));
    finishStep();
  }

  function onSkip() {
    finishStep();
  }

  return (
    <ScrollView className="flex-1 bg-white" contentContainerClassName="flex-1 justify-center px-6 py-12">
      <View className="w-full max-w-[440px] gap-6 self-center">
        <StepProgressBar currentStep={2} totalSteps={2} />
        <View className="gap-1">
          <Text className="text-center text-2xl font-bold text-text-primary">
            Invite your team
          </Text>
          <Text className="text-center text-sm text-text-secondary">
            Work is better together — invite a few teammates to get started.
          </Text>
        </View>
        <View className="gap-3">
          {emails.map((email, i) => (
            <TextField
              key={i}
              placeholder="teammate@company.com"
              autoComplete="off"
              autoCapitalize="none"
              keyboardType="email-address"
              value={email}
              onChangeText={(value) => updateEmail(i, value)}
              error={errors[i]}
            />
          ))}
          <Pressable onPress={addAnother} className="self-start">
            <Text className="text-sm text-brand">+ Add another</Text>
          </Pressable>
        </View>
        <View className="gap-3">
          <PrimaryButton label="Send invite" onPress={onSendInvite} />
          <Pressable onPress={onSkip}>
            <Text className="text-center text-sm text-text-secondary">Skip for now</Text>
          </Pressable>
        </View>
      </View>
    </ScrollView>
  );
}
