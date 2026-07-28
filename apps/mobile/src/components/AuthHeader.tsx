import { Pressable, Text, View } from 'react-native';

interface AuthHeaderProps {
  navLabel?: string;
  onNavPress?: () => void;
}

// Wordmark matches apps/web/src/components/layout/Logo.tsx. Unlike web's
// Logo, this doesn't link anywhere — there's no landing page in the mobile
// app for it to point back to.
export function AuthHeader({ navLabel, onNavPress }: AuthHeaderProps) {
  return (
    <View className="flex-row items-center justify-between border-b border-gray-200 bg-white px-6 py-4">
      <View className="flex-row items-center gap-2">
        <View className="h-6 w-6 rounded bg-brand" />
        <Text className="text-[15px] font-bold text-text-primary">FoundryHQ</Text>
      </View>
      {navLabel && onNavPress ? (
        <Pressable onPress={onNavPress}>
          <Text className="text-sm text-text-secondary">{navLabel}</Text>
        </Pressable>
      ) : null}
    </View>
  );
}
