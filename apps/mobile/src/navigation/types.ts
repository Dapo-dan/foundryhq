export type AuthStackParamList = {
  SignIn: undefined;
  SignUp: undefined;
  ForgotPassword: undefined;
  ResetPassword: { token: string };
};

export type OnboardingStackParamList = {
  Workspace: undefined;
  Invite: undefined;
};

export type MainTabParamList = {
  Dashboard: undefined;
};
