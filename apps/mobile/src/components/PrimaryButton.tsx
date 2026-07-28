import { ActivityIndicator, Pressable, Text, type PressableProps } from 'react-native';

export interface PrimaryButtonProps extends PressableProps {
  label: string;
  loading?: boolean;
}

export function PrimaryButton({ label, loading, disabled, ...pressableProps }: PrimaryButtonProps) {
  return (
    <Pressable
      className="h-11 items-center justify-center rounded-lg bg-brand-navy active:opacity-80 disabled:opacity-50"
      disabled={disabled || loading}
      {...pressableProps}
    >
      {loading ? (
        <ActivityIndicator color="#FFFFFF" />
      ) : (
        <Text className="text-base font-medium text-white">{label}</Text>
      )}
    </Pressable>
  );
}
