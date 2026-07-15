import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppLayout } from '@/components/layout'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { LandingPage } from '@/pages/landing'
import { PrivacyPage } from '@/pages/privacy'
import { TermsPage } from '@/pages/terms'
import { DashboardPage } from '@/pages/dashboard'
import { ProjectsPage } from '@/pages/projects'
import { TasksPage } from '@/pages/tasks'
import { SettingsPage } from '@/pages/settings'
import { SignInPage } from '@/pages/auth/sign-in'
import { SignUpPage } from '@/pages/auth/sign-up'
import { OnboardingPage } from '@/pages/onboarding'
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
      { path: 'settings', element: <SettingsPage /> },
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
        ],
      },
      { path: 'onboarding', element: <OnboardingPage /> },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
