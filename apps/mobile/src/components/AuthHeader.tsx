import { Pressable, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

interface AuthHeaderProps {
  navLabel?: string;
  onNavPress?: () => void;
}

// Wordmark matches apps/web/src/components/layout/Logo.tsx. Unlike web's
// Logo, this doesn't link anywhere — there's no landing page in the mobile
// app for it to point back to.
//
// With the native stack header hidden (see AuthNavigator), nothing else
// accounts for the status bar/notch — SafeAreaView (top edge only, since
// this is the top bar) keeps the wordmark from rendering under it.
export function AuthHeader({ navLabel, onNavPress }: AuthHeaderProps) {
  return (
    <SafeAreaView edges={['top']} className="border-b border-gray-200 bg-white">
      <View className="flex-row items-center justify-between px-6 py-4">
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
    </SafeAreaView>
  );
}
