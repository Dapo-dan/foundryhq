import type { ReactNode } from 'react';
import { Text, TextInput, View, type TextInputProps } from 'react-native';

export interface TextFieldProps extends TextInputProps {
  label?: string;
  error?: string;
  rightElement?: ReactNode;
}

export function TextField({ label, error, rightElement, ...inputProps }: TextFieldProps) {
  return (
    <View className="gap-1.5">
      {label ? <Text className="text-sm font-medium text-text-primary">{label}</Text> : null}
      <View className="relative">
        <TextInput
          className={`h-11 rounded-lg border px-3 text-base text-text-primary ${
            error ? 'border-red-600' : 'border-gray-300'
          } ${rightElement ? 'pr-14' : ''}`}
          placeholderTextColor="#888888"
          {...inputProps}
        />
        {rightElement ? (
          <View className="absolute right-3 top-0 h-11 items-center justify-center">
            {rightElement}
          </View>
        ) : null}
      </View>
      {error ? <Text className="text-sm text-red-600">{error}</Text> : null}
    </View>
  );
}
