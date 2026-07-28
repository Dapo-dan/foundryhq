import { useState } from 'react';
import { Pressable, Text } from 'react-native';
import { TextField, type TextFieldProps } from './TextField';

type PasswordFieldProps = Omit<TextFieldProps, 'secureTextEntry' | 'rightElement'>;

// Not from any component library — a small addition on top of TextField for
// the "Show/Hide" text toggle the web auth/onboarding forms use instead of
// an eye icon (see apps/web/src/components/ui/password-input.tsx).
export function PasswordField(props: PasswordFieldProps) {
  const [visible, setVisible] = useState(false);

  return (
    <TextField
      {...props}
      secureTextEntry={!visible}
      rightElement={
        <Pressable onPress={() => setVisible((v) => !v)}>
          <Text className="text-xs font-medium text-brand">{visible ? 'Hide' : 'Show'}</Text>
        </Pressable>
      }
    />
  );
}
