import { createBrowserRouter } from 'react-router-dom'
import { NotFoundPage } from '@/pages/not-found'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <div>FoundryHQ</div>,
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
