# 0001. Use a monorepo for API, web, and mobile

**Status:** Accepted
**Date:** 2026-07-10
**Deciders:** Founding engineering team

## Context

FoundryHQ ships three deployables — a Go API, a React web SPA, and a React Native mobile app — that share a data model and, in large part, a type system. Early on we had to decide whether each would live in its own repository or all three would live together.

## Decision

Use a single monorepo (`foundryhq/`) containing `apps/api`, `apps/web`, `apps/mobile`, and a shared `packages/shared-types` package, with one CI pipeline covering all three.

## Alternatives Considered

- **Separate repos per app** — rejected: shared types between web and mobile would require either a published npm package (slow iteration — every type change needs a publish/bump/install cycle) or manual duplication (guaranteed drift).
- **Separate repos with a shared-types repo** — rejected: three extra repos to keep in sync for a team this size adds process overhead disproportionate to the benefit; revisit if/when teams grow large enough to need independent release cadences and access control per app.

## Consequences

- Cross-app refactors (e.g., a shared type change) land in a single PR instead of a multi-repo coordination dance
- `packages/shared-types` is trivially kept in sync since it's just another workspace package, not a published dependency
- CI must be scoped carefully (path-based triggers) so a web-only change doesn't force a full API test suite run
- As the team grows, per-app access control and independent versioning (see `../release.md`) become harder inside one repo — worth revisiting if/when engineering headcount grows significantly past founding-team scale

## Related

- Architecture: `../architecture.md#monorepo-layout`
- Release process: `../../.ai/documentation/release.md`
