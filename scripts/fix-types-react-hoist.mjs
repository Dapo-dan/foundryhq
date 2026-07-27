// apps/web (React 18) and apps/mobile (React 19) each resolve their own
// @types/react version correctly, but lucide-react has no @types/react of
// its own (react is a peer, @types/react isn't declared at all) — so
// TypeScript falls back to pnpm's shared `.pnpm/node_modules` peer-resolution
// slot when checking lucide-react's JSX usage in apps/web. That shared slot
// only holds one version workspace-wide, and pnpm doesn't guarantee it picks
// web's — it has landed on mobile's React 19 types before, which breaks
// web's build (React 19's ReactNode allows `bigint`, which React 18's
// doesn't, so every lucide-react icon fails JSX type-checking).
//
// No pnpm config (packageExtensions, overrides, resolve-peers-from-workspace-root)
// or tsconfig setting (typeRoots) changes which version wins that slot — this
// script runs after every install and repoints it to whatever apps/web
// itself actually resolves, which is the one thing confirmed to fix it.

import { lstatSync, symlinkSync, unlinkSync, realpathSync } from 'node:fs'
import { dirname, join, relative } from 'node:path'
import { fileURLToPath } from 'node:url'

const repoRoot = dirname(dirname(fileURLToPath(import.meta.url)))
const webTypesReact = join(repoRoot, 'apps/web/node_modules/@types/react')
const fallbackTypesReact = join(repoRoot, 'node_modules/.pnpm/node_modules/@types/react')

function tryStat(path) {
  try {
    return lstatSync(path)
  } catch {
    return null
  }
}

if (!tryStat(webTypesReact) || !tryStat(fallbackTypesReact)) {
  // Nothing to fix — either apps/web has no @types/react resolved yet, or
  // this pnpm version doesn't create the shared fallback slot at all.
  process.exit(0)
}

const correctTarget = realpathSync(webTypesReact)
const currentTarget = tryStat(fallbackTypesReact).isSymbolicLink()
  ? realpathSync(fallbackTypesReact)
  : null

if (currentTarget === correctTarget) {
  process.exit(0)
}

unlinkSync(fallbackTypesReact)
symlinkSync(relative(dirname(fallbackTypesReact), correctTarget), fallbackTypesReact)

console.log(
  `[fix-types-react-hoist] repointed the shared @types/react fallback to match apps/web (was pointing elsewhere)`
)
