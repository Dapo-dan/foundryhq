import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { Text, View } from 'react-native';
import type { AuthStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<AuthStackParamList, 'ResetPassword'>;

export function ResetPasswordScreen({ route }: Props) {
  return (
    <View className="flex-1 items-center justify-center bg-white">
      <Text className="text-text-primary">Reset Password</Text>
      <Text className="text-text-muted">token: {route.params.token}</Text>
    </View>
  );
}
