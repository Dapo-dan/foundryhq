import type { ReactNode } from 'react';
import { Text, View } from 'react-native';

interface AuthCardProps {
  heading: string;
  description: string;
  children: ReactNode;
}

// Shared shell for the 4 auth screens (Sign In, Sign Up, Forgot/Reset
// Password) — matches apps/web/src/components/layout/AuthCard.tsx's plain
// centered heading + subtext + form stack, no card border.
export function AuthCard({ heading, description, children }: AuthCardProps) {
  return (
    <View className="gap-6">
      <View className="gap-1">
        <Text className="text-center text-2xl font-bold text-text-primary">{heading}</Text>
        <Text className="text-center text-sm text-text-secondary">{description}</Text>
      </View>
      {children}
    </View>
  );
}
