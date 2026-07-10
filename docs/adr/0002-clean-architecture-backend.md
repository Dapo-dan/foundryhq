# 0002. Clean Architecture for the Go backend

**Status:** Accepted
**Date:** 2026-07-10
**Deciders:** Founding engineering team

## Context

The API needs business logic (deal stage transitions, sprint velocity calculation, OKR roll-ups) that must be testable independent of HTTP and the database, and must not leak infrastructure details (Gin request objects, GORM models) into that logic. Without an enforced boundary, business rules tend to accrete inside handlers where they're hard to unit test and hard to reuse.

## Decision

Structure `apps/api` as Clean Architecture with a strict inward dependency direction: `handlers → usecases → domain ← repositories`. `domain` defines entities and repository interfaces; `usecases` contain all business logic and depend only on domain interfaces; `repositories` implement those interfaces against PostgreSQL; `handlers` bind HTTP requests and call usecases.

## Alternatives Considered

- **Fat handlers, thin everything else** — rejected: fastest to write initially, but business logic becomes untestable without spinning up `httptest` and a real database for every test, and logic reuse across handlers requires copy-paste.
- **MVC-style with "service" layer only (no domain/repository split)** — rejected: services calling GORM directly means every business-logic test needs a real database, and swapping the ORM or adding a cache later means rewriting business logic, not just the data layer.

## Consequences

- Usecases are fully unit-testable with mocked repositories, no database required (see `../../.ai/engineering/testing.md`)
- Every new feature requires touching four layers (domain, usecase, repository, handler) instead of one — more files for simple CRUD, deliberately so
- The dependency-direction rule is enforced by code review, not the compiler — a `handlers`-imports-into-`usecases` violation is a review-time catch, not a build-time one

## Related

- Architecture: `../architecture.md#backend-clean-architecture`
- Repository pattern: `0003-repository-pattern.md`
