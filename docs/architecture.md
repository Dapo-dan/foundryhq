# Architecture

Standalone reference for how FoundryHQ is built. For the day-to-day coding rules derived from this architecture, see `../.ai/engineering/architecture.md`; for why specific decisions were made, see `adr/`.

## Monorepo Layout

```
foundryhq/
├── apps/
│   ├── api/        # Go + Gin REST API
│   ├── web/        # React + TypeScript SPA
│   └── mobile/     # React Native + Expo
├── packages/
│   └── shared-types/   # TypeScript types shared by web & mobile
└── docs/                # this documentation set
```

See `adr/0001-monorepo-structure.md` for why this is one repository rather than three.

## Backend: Clean Architecture

Dependency direction is always **inward**:

```
handlers → usecases → domain ← repositories
```

- **`domain/`** — pure Go structs and repository interfaces. No HTTP, no DB, no imports from other internal packages.
- **`usecases/`** — business logic. Depends only on domain interfaces. Fully testable without a database or HTTP stack.
- **`repositories/postgres/`** — GORM implementations of domain interfaces. No business logic.
- **`handlers/`** — Gin HTTP handlers. Bind requests, call usecases, return JSON responses.
- **`middleware/`** — auth, CORS, logging, request ID.
- **`pkg/`** — reusable packages (JWT, database connection, logger) with no domain types.

See `adr/0002-clean-architecture-backend.md` and `adr/0003-repository-pattern.md` for the reasoning.

## Frontend: Web

- **State:** Zustand for UI/client state, TanStack Query for server state and caching (see `adr/0005-optimistic-ui-tanstack-query.md`)
- **API calls:** centralized in `services/` — components never call `fetch` directly
- **Components:** `components/ui/` for shadcn primitives, `components/` for product-specific composition
- **Pages:** thin, route-level orchestration in `pages/` — no business logic
- **Hooks:** reusable stateful logic in `hooks/`

## Frontend: Mobile

- Mirrors the web structure (screens ≈ pages, same hooks/services pattern)
- Shares types from `packages/shared-types/` — types are never duplicated between web and mobile
- Navigation: React Navigation (stack + tab navigators)
- State: same Zustand + TanStack Query pattern as web

## Auth Flow

- Access token: JWT, 15-minute expiry, held in memory (web) / SecureStore (mobile)
- Refresh token: 7-day expiry, httpOnly cookie (web) / SecureStore (mobile)
- All protected routes require `Authorization: Bearer <token>`
- Auth middleware lives in `apps/api/internal/middleware/`

See `adr/0004-jwt-access-refresh-tokens.md` for why this shape was chosen over server-side sessions.

## Multi-Tenancy

Every workspace-scoped table carries a `workspace_id` foreign key, and every query is filtered by the authenticated user's workspace membership at the repository layer — never left to the handler or the client to enforce. See `database.md` for the schema this produces.

## Data Flow

```
Browser / Mobile App
        │
        ▼
   REST API (Gin)
        │
   Middleware (auth, CORS, logging, request ID)
        │
        ▼
   Handler (bind & validate request)
        │
        ▼
   Usecase (business logic)
        │
        ▼
   Repository Interface
        │
        ▼
   PostgreSQL (via GORM)
```

## Key Architectural Rules

1. Never import from `handlers` into `usecases` or `domain` — the dependency arrow never points outward.
2. Repository interfaces are defined in `domain`, not in `repositories` — the domain owns the contract.
3. No raw SQL — use GORM. If GORM can't express a query cleanly, write a raw query in the repository layer with a comment explaining why.
4. One usecase per business operation — fat usecases are fine; fat handlers are not.
5. Shared TypeScript types come from `packages/shared-types` — never copy-paste type definitions between web and mobile.

## Adding a New Feature (Full-Stack Checklist)

- [ ] Domain struct + repository interface in `domain/`
- [ ] Migration + GORM model in `repositories/postgres/` (see `database.md` for schema conventions)
- [ ] Usecase with business logic in `usecases/`
- [ ] Gin handler + routes in `handlers/`, documented per `../.ai/documentation/api.md`
- [ ] TypeScript type in `packages/shared-types/`
- [ ] Service function + TanStack Query hooks in `apps/web/src/`
- [ ] Web page/component
- [ ] Mobile screen in `apps/mobile/src/screens/`
- [ ] Requirement traced to a REQ-ID in `requirements.md`
