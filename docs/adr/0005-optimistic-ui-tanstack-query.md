# 0005. Optimistic UI with TanStack Query

**Status:** Accepted
**Date:** 2026-07-10
**Deciders:** Web engineering

## Context

FoundryHQ's core interactions — moving a deal between pipeline stages, checking off a task, updating an OKR check-in — happen constantly during a founder's day. Waiting for a server round-trip before the UI reflects a change makes the product feel slow for exactly the interactions that need to feel instant (see the "speed over completeness" principle in `../vision.md`).

## Decision

Use TanStack Query for all server state on web and mobile. Mutations update the local query cache optimistically before the server responds, and roll back to the previous cache state with a toast notification if the request fails.

## Alternatives Considered

- **Wait for server confirmation before updating UI** — rejected: correct and simple, but every drag-and-drop or status toggle would show a loading flicker, directly against the product's speed principle.
- **Hand-rolled optimistic state in Zustand** — rejected: TanStack Query already solves cache invalidation, retry, and rollback; reimplementing it in the global store duplicates a solved problem and would need to be kept in sync by hand across every mutation.

## Consequences

- Every mutation hook must define an `onError` rollback — a mutation added without one will show a stale success state that silently reverts on refresh, which is worse than no optimism at all
- Race conditions between rapid successive mutations on the same entity (e.g., fast drag-and-drop across stages) need explicit handling via TanStack Query's mutation queuing/cancellation, not left implicit
- Server responses remain the source of truth on refetch — optimistic state is a UI convenience, never persisted as if confirmed

## Related

- Architecture: `../architecture.md#frontend-web`
- Requirements: `REQ-NFR-04` in `../requirements.md`
