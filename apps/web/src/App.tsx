import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from 'react-router-dom'
import { router } from '@/routes'

const queryClient = new QueryClient()

// A React "component" is just a function that returns JSX (the HTML-like
// syntax below). React calls this function to figure out what to display,
// and re-calls it whenever the component's state/props change. Component
// function names must start with a capital letter — that's how React (and
// JSX) tells your own components apart from plain HTML tags like <div>.
export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  )
}

// main.tsx imports this as `App` and renders it — this is the root/top-level
// component that everything else in the app will eventually nest inside.
