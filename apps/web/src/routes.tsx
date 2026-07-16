import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppLayout } from '@/components/layout'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { OnboardingLayout } from '@/components/layout/OnboardingLayout'
import { LandingPage } from '@/pages/landing'
import { PrivacyPage } from '@/pages/privacy'
import { TermsPage } from '@/pages/terms'
import { DashboardPage } from '@/pages/dashboard'
import { ProjectsPage } from '@/pages/projects'
import { TasksPage } from '@/pages/tasks'
import { SettingsPage } from '@/pages/settings'
import { DealPipelinePage } from '@/pages/crm/deal-pipeline'
import { CompaniesPage } from '@/pages/crm/companies'
import { ContactsPage } from '@/pages/crm/contacts'
import { SprintsPage } from '@/pages/sprints'
import { MetricsPage } from '@/pages/kpis/metrics'
import { ReportsPage } from '@/pages/kpis/reports'
import { SignInPage } from '@/pages/auth/sign-in'
import { SignUpPage } from '@/pages/auth/sign-up'
import { ForgotPasswordPage } from '@/pages/auth/forgot-password'
import { ResetPasswordPage } from '@/pages/auth/reset-password'
import { WorkspaceStepPage } from '@/pages/onboarding/workspace'
import { TeamSizeStepPage } from '@/pages/onboarding/team-size'
import { RoleStepPage } from '@/pages/onboarding/role'
import { ToolsStepPage } from '@/pages/onboarding/tools'
import { InviteStepPage } from '@/pages/onboarding/invite'
import { WelcomeStepPage } from '@/pages/onboarding/welcome'
import { NotFoundPage } from '@/pages/not-found'

export const router = createBrowserRouter([
  { path: '/', element: <LandingPage /> },
  { path: 'privacy', element: <PrivacyPage /> },
  { path: 'terms', element: <TermsPage /> },
  {
    element: <AppLayout />,
    children: [
      { path: 'dashboard', element: <DashboardPage /> },
      { path: 'projects', element: <ProjectsPage /> },
      { path: 'tasks', element: <TasksPage /> },
      { path: 'sprints', element: <SprintsPage /> },
      { path: 'settings', element: <SettingsPage /> },
      { path: 'crm/deal-pipeline', element: <DealPipelinePage /> },
      { path: 'crm/companies', element: <CompaniesPage /> },
      { path: 'crm/contacts', element: <ContactsPage /> },
      { path: 'kpis/metrics', element: <MetricsPage /> },
      { path: 'kpis/reports', element: <ReportsPage /> },
    ],
  },
  {
    element: <AuthLayout />,
    children: [
      {
        path: 'auth',
        children: [
          { index: true, element: <Navigate to="sign-in" replace /> },
          { path: 'sign-in', element: <SignInPage /> },
          { path: 'sign-up', element: <SignUpPage /> },
          { path: 'forgot-password', element: <ForgotPasswordPage /> },
          { path: 'reset-password', element: <ResetPasswordPage /> },
        ],
      },
      {
        path: 'onboarding',
        children: [
          { index: true, element: <Navigate to="workspace" replace /> },
          {
            element: <OnboardingLayout />,
            children: [
              { path: 'workspace', element: <WorkspaceStepPage /> },
              { path: 'team-size', element: <TeamSizeStepPage /> },
              { path: 'role', element: <RoleStepPage /> },
              { path: 'tools', element: <ToolsStepPage /> },
              { path: 'invite', element: <InviteStepPage /> },
              { path: 'welcome', element: <WelcomeStepPage /> },
            ],
          },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
