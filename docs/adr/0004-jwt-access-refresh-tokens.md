# 0004. JWT access + refresh tokens for authentication

**Status:** Accepted
**Date:** 2026-07-10
**Deciders:** Founding engineering team

## Context

FoundryHQ needs stateless auth that works across a web SPA and a React Native mobile app, without requiring server-side session storage that complicates horizontally scaling the Go API.

## Decision

Use JWT access tokens (15-minute expiry) plus refresh tokens (7-day expiry). Access tokens are held in memory on web and in SecureStore on mobile; refresh tokens are stored in an httpOnly cookie on web and SecureStore on mobile.

## Alternatives Considered

- **Server-side sessions (Redis-backed)** — rejected: adds a stateful dependency and a network hop to every authenticated request, and complicates mobile token storage since there's no cookie-jar equivalent.
- **Long-lived single JWT, no refresh** — rejected: a stolen long-lived token can't be revoked short of a blocklist, which reintroduces server-side state anyway, without the security benefit of short-lived access tokens.

## Consequences

- Logout requires explicit refresh-token invalidation server-side — a small stateful surface remains, but far smaller than full sessions
- Every client needs correct silent-refresh handling — a missed refresh surfaces as a hard logout, not a soft error
- Web and mobile share the same token contract, keeping the auth types in `packages/shared-types/` identical across platforms

## Related

- Architecture: `../architecture.md#auth-flow`
- Requirements: `REQ-01`, `REQ-NFR-03` in `../requirements.md`
