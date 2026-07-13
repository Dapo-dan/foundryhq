import { createBrowserRouter } from 'react-router-dom'
import { AppLayout } from '@/components/layout'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { LandingPage } from '@/pages/landing'
import { DashboardPage } from '@/pages/dashboard'
import { ProjectsPage } from '@/pages/projects'
import { TasksPage } from '@/pages/tasks'
import { SettingsPage } from '@/pages/settings'
import { AuthPage } from '@/pages/auth'
import { OnboardingPage } from '@/pages/onboarding'
import { NotFoundPage } from '@/pages/not-found'

export const router = createBrowserRouter([
  { path: '/', element: <LandingPage /> },
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
      { path: 'auth', element: <AuthPage /> },
      { path: 'onboarding', element: <OnboardingPage /> },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
