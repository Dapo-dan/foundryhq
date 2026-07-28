import { getPasswordStrengthScore } from '@foundryhq/shared-validation';
import { Text, View } from 'react-native';

const LABELS = ['Weak', 'Fair', 'Good', 'Strong'];

interface PasswordStrengthBarProps {
  password: string;
}

export function PasswordStrengthBar({ password }: PasswordStrengthBarProps) {
  const score = getPasswordStrengthScore(password);

  if (!password) return null;

  return (
    <View className="flex-row items-center gap-2">
      <View className="flex-1 flex-row gap-1">
        {Array.from({ length: 4 }, (_, i) => (
          <View
            key={i}
            className={`h-1 flex-1 rounded-full ${i < score ? 'bg-brand-accent' : 'bg-gray-200'}`}
          />
        ))}
      </View>
      <Text className="text-xs text-text-subtle">{LABELS[Math.max(score - 1, 0)]}</Text>
    </View>
  );
}
