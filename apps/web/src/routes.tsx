import { createBrowserRouter } from 'react-router-dom'
import { AppLayout } from '@/components/layout'
import { NotFoundPage } from '@/pages/not-found'

export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    children: [{ index: true, element: <div>Dashboard placeholder</div> }],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
