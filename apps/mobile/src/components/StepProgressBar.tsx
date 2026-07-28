import { View } from 'react-native';

interface StepProgressBarProps {
  /** 1-indexed */
  currentStep: number;
  totalSteps: number;
}

export function StepProgressBar({ currentStep, totalSteps }: StepProgressBarProps) {
  return (
    <View className="flex-row gap-1.5">
      {Array.from({ length: totalSteps }, (_, i) => (
        <View
          key={i}
          className={`h-1 flex-1 rounded-full ${i < currentStep ? 'bg-brand-navy' : 'bg-gray-200'}`}
        />
      ))}
    </View>
  );
}
