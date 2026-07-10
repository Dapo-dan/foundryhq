# 0003. Repository pattern for data access

**Status:** Accepted
**Date:** 2026-07-10
**Deciders:** Founding engineering team

## Context

Given the Clean Architecture decision (`0002-clean-architecture-backend.md`), usecases need a way to persist and query data without depending on GORM or PostgreSQL directly — otherwise "testable without a database" is not actually achievable.

## Decision

Define repository interfaces in `domain/` (e.g., `TaskRepository`, `DealRepository`), owned by the domain layer, and implement them in `repositories/postgres/` using GORM. Usecases depend only on the interface type.

## Alternatives Considered

- **Usecases call GORM directly** — rejected: ties every business-logic test to a running Postgres instance, and any future change to the ORM or database ripples into business logic instead of staying contained to the data layer.
- **Generic repository (single `Repository[T]` interface for all entities)** — rejected: FoundryHQ's queries are domain-specific enough (e.g., `GetByWorkspace` with filters, `GetOverdueForFollowUp`) that a generic CRUD interface would just get bypassed with raw queries anyway, defeating the point.

## Consequences

- Mock repositories for unit tests are simple hand-written structs implementing the domain interface — no test-database setup needed for usecase tests
- Swapping PostgreSQL for another store (unlikely, but e.g. adding a read replica or a cache layer) touches only `repositories/postgres/`, not business logic
- Slight duplication: each entity needs both a domain interface and a GORM struct in the repository layer, kept in sync by convention rather than a shared type — deviation here (interface says one thing, GORM implementation does another) is a code-review catch, not a compiler error

## Related

- Architecture: `../architecture.md#backend-clean-architecture`
- Database schema: `../database.md`
