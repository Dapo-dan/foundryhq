import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from 'react-router-dom'
import { Toaster } from 'sonner'
import { router } from '@/routes'
import { useSessionBootstrap } from '@/hooks/useSessionBootstrap'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
})

// A React "component" is just a function that returns JSX (the HTML-like
// syntax below). React calls this function to figure out what to display,
// and re-calls it whenever the component's state/props change. Component
// function names must start with a capital letter — that's how React (and
// JSX) tells your own components apart from plain HTML tags like <div>.
export function App() {
  useSessionBootstrap()

  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <Toaster richColors position="bottom-right" />
    </QueryClientProvider>
  )
}

// main.tsx imports this as `App` and renders it — this is the root/top-level
// component that everything else in the app will eventually nest inside.
