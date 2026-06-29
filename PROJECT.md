# FoundryHQ — Project Overview

## 1. Project Vision

Most startups cobble together 10+ disconnected tools — a CRM here, a task tracker there, a notes app somewhere else. Context gets lost, time is wasted, and teams lose alignment as they scale.

**FoundryHQ** is the unified operating system built specifically for startups. One platform where your team tracks customers, ships product, runs meetings, and measures what matters — from day zero to Series B.

The goal is radical simplicity: everything a founding team needs, nothing they don't, in a single fast workspace that grows with them.

---

## 2. Target Users

| Persona | Pain Point | How FoundryHQ Helps |
|---|---|---|
| **Founding teams (2–10 people)** | Too many tools, not enough time | Single workspace replaces Notion + Linear + HubSpot |
| **Startup operators / COOs** | No single view of company health | Real-time KPI dashboard tied to actual work |
| **Early sales teams** | CRM overkill (Salesforce) or underbuilt (spreadsheets) | Lightweight CRM designed for founder-led sales |
| **Product-focused engineers** | Context switching between PM tools and code | Sprint tracking that connects to real team goals |
| **Investors / advisors** | No transparent view into portfolio company ops | Shareable KPI and goal progress reports |

---

## 3. Features

### Authentication
- Email + password registration and login
- OAuth 2.0 (Google, GitHub)
- JWT access and refresh tokens
- Role-based access: Owner, Admin, Member, Viewer

### CRM
- Contact and company management
- Deal pipeline with Kanban view and drag-and-drop stages
- Activity timeline (calls, emails, notes, meetings)

### Task Management
- Kanban board view (by status)
- Sprint planning — create sprints, assign tasks, track velocity
- Sub-tasks, labels, priorities (Urgent / High / Medium / Low), due dates
- Assignees and watchers

### Meeting Notes
- Create structured meeting records with attendees
- Rich-text notes
- Action items with assignees and due dates
- Link meetings to projects or CRM contacts
- Searchable archive

### Goal Tracking (OKRs)
- Company, team, and personal OKRs
- Key results with numeric, percentage, or boolean metrics
- Progress auto-calculated from key result check-ins
- Alignment tree — see how team goals roll up to company goals

### Team Management
- Invite members via email
- Role management per workspace
- Team activity feed
- Workspace settings (name, logo, slug)

### KPI Dashboard
- Define custom KPIs with targets (number, percent, currency)
- Record values over time (time-series)
- Link KPIs to goals for alignment visibility
- Progress vs. target visualisation

### Notifications
- In-app notifications
- Push notifications (mobile)
- Email notifications
- @mention support across tasks, meetings, and goals

---

## 4. Architecture

### Monorepo Structure

```
foundryhq/
├── apps/
│   ├── api/                        # Go + Gin REST API
│   │   ├── cmd/
│   │   │   └── server/             # main.go — entrypoint
│   │   ├── internal/
│   │   │   ├── domain/             # Core entities & repository interfaces
│   │   │   ├── usecases/           # Business logic (no HTTP, no DB details)
│   │   │   ├── repositories/       # PostgreSQL implementations
│   │   │   │   └── postgres/
│   │   │   └── handlers/           # Gin HTTP handlers
│   │   ├── internal/
│   │   │   ├── middleware/         # Auth, CORS, logging, request ID
│   │   │   └── config/             # Environment config loader
│   │   └── pkg/                    # Reusable packages (jwt, database, logger)
│   ├── web/                        # React + TypeScript SPA
│   │   └── src/
│   │       ├── components/         # Reusable UI components
│   │       ├── pages/              # Route-level page components
│   │       ├── hooks/              # Custom React hooks
│   │       ├── store/              # Zustand global state slices
│   │       ├── services/           # API client functions
│   │       ├── types/              # TypeScript type definitions
│   │       └── lib/                # Utilities
│   └── mobile/                     # React Native + Expo
│       └── src/
│           ├── components/
│           ├── screens/
│           ├── navigation/
│           ├── services/
│           ├── store/
│           ├── types/
│           └── hooks/
├── packages/
│   └── shared-types/               # TypeScript types shared by web & mobile
├── docs/
│   ├── api/                        # OpenAPI / Swagger spec
│   ├── architecture/               # Architecture decision records (ADRs)
│   └── guides/                     # Developer guides
└── .github/
    └── workflows/                  # CI/CD pipelines (lint, test, build, deploy)
```

### Backend — Clean Architecture

The API follows Clean Architecture. Dependency direction is always inward:

```
handlers → usecases → domain ← repositories
```

- **domain** — pure Go structs and repository interfaces; knows nothing about HTTP or databases
- **usecases** — all business logic; depends only on domain interfaces
- **repositories** — PostgreSQL implementations of domain interfaces; no business logic
- **handlers** — Gin HTTP handlers; bind requests, call usecases, return responses

### Tech Stack

| Layer | Technology | Reason |
|---|---|---|
| **Frontend** | React 18, TypeScript, Vite, Tailwind CSS, shadcn/ui | Fast DX, type safety, utility-first styling |
| **State (web)** | Zustand + TanStack Query | Zustand for UI state; TQ for server state & caching |
| **Mobile** | React Native, Expo, React Navigation | Code reuse with web types, fast iteration |
| **Backend** | Go 1.22+, Gin | Performance, simplicity, strong concurrency |
| **ORM / DB** | GORM, PostgreSQL 16 | Mature, reliable, relational fits the domain |
| **Auth** | JWT (access + refresh), OAuth 2.0 | Stateless, standard, mobile-compatible |
| **API Docs** | Swagger / OpenAPI 3.0 (swaggo) | Auto-generated from Go annotations |
| **Infra** | Docker, Docker Compose | Reproducible local environment |
| **CI/CD** | GitHub Actions | Lint, test, build, and deploy on push |

### Data Flow

```
Browser / Mobile App
        │
        ▼
   REST API (Gin)
        │
   Middleware (auth, cors, logging)
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

### Key Design Decisions

- **Monorepo** — shared types between web and mobile; unified CI pipeline; easier cross-app refactors
- **Clean Architecture on backend** — usecases are fully testable without a database or HTTP stack
- **Repository pattern** — swap the DB implementation without touching business logic
- **JWT with refresh tokens** — short-lived access tokens (15 min) + long-lived refresh tokens (7 days) stored in httpOnly cookies on web
- **Optimistic UI** — TanStack Query mutations update the UI immediately and roll back on error
