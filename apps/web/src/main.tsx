import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { App } from './App.tsx'

// This is the entry point of the whole app — the one file every React app
// needs. It finds the empty <div id="root"> in index.html and tells React
// "render our <App /> component tree inside this DOM element." Everything
// you build (pages, components, routes) will hang off of <App />.
createRoot(document.getElementById('root')!).render(
  // StrictMode doesn't render anything itself — it's a dev-only wrapper that
  // makes React double-invoke some functions (like component bodies and
  // effects) on purpose, to help you catch bugs where code isn't "pure"
  // (e.g. relies on running only once). It has no effect in production.
  <StrictMode>
    <App />
  </StrictMode>,
)
