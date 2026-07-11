import path from 'node:path'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// Vite is the dev server + bundler (not a React-specific thing — this is
// "build tooling," separate from React itself). It serves your app locally
// with instant hot-reload, and bundles everything into static files for
// `pnpm build`.
// https://vite.dev/config/
export default defineConfig({
  // Plugins teach Vite how to handle things it doesn't understand natively:
  // - react(): compiles JSX syntax (the <div>...</div> in .tsx files) into
  //   plain JS function calls, and enables fast-refresh (edit a component,
  //   see the change instantly without losing app state).
  // - tailwindcss(): scans your className strings and generates only the CSS
  //   for the utility classes you actually used.
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      // Lets you write `import { cn } from '@/lib/utils'` instead of a
      // relative path like '../../lib/utils' — '@' always points at src/.
      '@': path.resolve(__dirname, './src'),
    },
  },
})
