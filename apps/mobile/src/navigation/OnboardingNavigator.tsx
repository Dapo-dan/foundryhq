import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { InviteScreen } from '../screens/onboarding/InviteScreen';
import { WorkspaceScreen } from '../screens/onboarding/WorkspaceScreen';
import { useOnboardingStore } from '../store/slices/onboarding';
import type { OnboardingStackParamList } from './types';

const Stack = createNativeStackNavigator<OnboardingStackParamList>();

export function OnboardingNavigator() {
  // Resume at the first incomplete step rather than always restarting at
  // Workspace — e.g. a user who named their workspace but backgrounded the
  // app before inviting anyone should land back on Invite, not repeat Workspace.
  const isWorkspaceDone = useOnboardingStore((state) => state.completedSteps.includes('workspace'));

  return (
    // Each screen already renders its own heading text — the native stack
    // header would just duplicate it with a second, differently-worded title.
    <Stack.Navigator
      initialRouteName={isWorkspaceDone ? 'Invite' : 'Workspace'}
      screenOptions={{ headerShown: false }}
    >
      <Stack.Screen name="Workspace" component={WorkspaceScreen} />
      <Stack.Screen name="Invite" component={InviteScreen} />
    </Stack.Navigator>
  );
}
