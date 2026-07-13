import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppLayout } from '@/components/layout'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { DashboardPage } from '@/pages/dashboard'
import { ProjectsPage } from '@/pages/projects'
import { TasksPage } from '@/pages/tasks'
import { SettingsPage } from '@/pages/settings'
import { AuthPage } from '@/pages/auth'
import { NotFoundPage } from '@/pages/not-found'

export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <DashboardPage /> },
      { path: 'projects', element: <ProjectsPage /> },
      { path: 'tasks', element: <TasksPage /> },
      { path: 'settings', element: <SettingsPage /> },
    ],
  },
  {
    element: <AuthLayout />,
    children: [{ path: 'auth', element: <AuthPage /> }],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
