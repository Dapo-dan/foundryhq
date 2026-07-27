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
    <Stack.Navigator initialRouteName={isWorkspaceDone ? 'Invite' : 'Workspace'}>
      <Stack.Screen
        name="Workspace"
        component={WorkspaceScreen}
        options={{ title: 'Create a Workspace' }}
      />
      <Stack.Screen
        name="Invite"
        component={InviteScreen}
        options={{ title: 'Invite Your Team' }}
      />
    </Stack.Navigator>
  );
}
