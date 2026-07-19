# FoundryHQ

> The complete operating system for modern startups.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![CI](https://github.com/foundryhq/foundryhq/actions/workflows/ci.yml/badge.svg)](https://github.com/foundryhq/foundryhq/actions)

---

## Status

Auth is built end-to-end (API + web). The web app also has UI for CRM, tasks, goals, meetings, team, dashboard, and KPIs, but those currently run on mock data (`apps/web/src/lib/mock/`) rather than a live backend — only auth is wired to the real API so far. `apps/mobile` has no code yet, just directory scaffolding. The **Features** section below describes the target product; see [`docs/mvp.md`](docs/mvp.md) for what's actually being built first (Auth, Workspace/Team, Tasks) and [`docs/roadmap.md`](docs/roadmap.md) for sequencing after that.

---

## Vision

Most startups cobble together 10+ disconnected tools — a CRM here, a task tracker there, a notes app somewhere else. Context gets lost, time is wasted, and teams lose alignment as they scale.

**FoundryHQ** is the unified operating system built specifically for startups. One platform where your team tracks customers, ships product, runs meetings, and measures what matters — from day zero to Series B.

---

## Target Users

| Persona | Pain Point | How FoundryHQ Helps |
|---|---|---|
| **Founder-Operator** | Context scattered across five+ tools; no single view of company health | One workspace for CRM, tasks, goals, and KPIs — replaces Notion + Linear + HubSpot |
| **Product Engineer** | Sprint context lives in Slack/Linear, disconnected from real goals | Sprint tracking tied to actual OKRs, keyboard-first workflow |
| **Early Sales Rep** | CRM overkill (Salesforce) or underbuilt (spreadsheets) | Lightweight deal pipeline designed for founder-led sales |
| **Investor/Advisor** | No transparent, read-only view into portfolio company ops | Shareable KPI/OKR snapshot — no full workspace login required |

---

## Features

### Authentication
- Email + password registration and login
- OAuth 2.0 (Google, GitHub)
- JWT access and refresh tokens
- Role-based access: Owner, Admin, Member, Viewer

### CRM
- Contact and company management
- Deal pipeline with Kanban view
- Activity timeline (calls, emails, notes)
- Integration-ready (email sync, calendar)

### Task Management
- Kanban and list views
- Sprint planning and backlog
- Sub-tasks, labels, priorities, due dates
- Assignees and watchers

### Meeting Notes
- Structured note templates
- Automatic action item extraction
- Link meetings to projects and contacts
- Searchable archive

### Goal Tracking (OKRs)
- Company, team, and personal OKRs
- Key result check-ins
- Progress visualization
- Alignment tree view

### Team Management
- Role management per workspace (see Authentication)
- Invite via email or link
- Team activity feed
- Onboarding checklists

### KPI Dashboard
- Configurable metric widgets
- Goal vs. actual tracking
- Time-series charts
- Export to PDF / CSV

### Notifications
- In-app, push (mobile), and email
- Digest preferences
- @mention support

---

## Architecture

```
foundryhq/                          # Monorepo root
├── apps/
│   ├── api/                        # Go + Gin REST API
│   │   ├── cmd/server/             # Entrypoint
│   │   ├── internal/
│   │   │   ├── domain/             # Core entities & interfaces
│   │   │   ├── usecases/           # Business logic
│   │   │   ├── repositories/       # Data access layer
│   │   │   └── handlers/           # HTTP handlers
│   │   └── pkg/                    # Reusable packages
│   ├── web/                        # React + TypeScript SPA
│   └── mobile/                     # React Native (iOS + Android)
├── packages/
│   └── shared-types/               # Shared TS types across web & mobile
├── docs/                           # Vision, requirements, architecture, roadmap, API, database, ADRs
├── docker-compose.yml
└── .github/workflows/              # CI/CD pipelines
```

### Tech Stack

| Layer | Technology |
|---|---|
| **Frontend** | React 18, TypeScript 5, Vite, Tailwind CSS, shadcn/ui, Zustand, React Query |
| **Mobile** | React Native 0.74, Expo, React Navigation |
| **Backend** | Go 1.22, Gin, GORM |
| **Database** | PostgreSQL 16 |
| **Auth** | JWT (access + refresh tokens), OAuth 2.0 (Google, GitHub) |
| **API Docs** | Swagger / OpenAPI 3.0 |
| **Infrastructure** | Docker, Docker Compose |
| **CI/CD** | GitHub Actions |

---

## Getting Started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [pnpm 9+](https://pnpm.io/installation)

### Quick Start (Docker)

`apps/api` and `apps/web` have no code yet (see **Status** above), so `docker-compose.yml` currently only brings up the database infra. The API and web services will be added to it once they're scaffolded.

```bash
git clone https://github.com/foundryhq/foundryhq.git
cd foundryhq

# Copy environment files
cp apps/api/.env.example apps/api/.env

# Start the database
docker compose up
```

- Postgres: localhost:5432
- pgAdmin: http://localhost:5050

Once `apps/api` and `apps/web` are scaffolded (see `docs/mvp.md`), this section will be updated with the commands to run them — targeting:

- Web app: http://localhost:5173
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

### Local Development

Not yet available — `apps/api`, `apps/web`, and `apps/mobile` have no code to run. This section will be filled in as each app is scaffolded per `docs/mvp.md`.

---

## API Documentation

The REST API is documented with Swagger/OpenAPI 3.0. See [`docs/api.md`](docs/api.md) for the curated endpoint reference.

- **Dev:** http://localhost:8080/swagger/index.html
- **Generated spec:** `docs/swagger/swagger.yaml`

Generate/update Swagger docs:
```bash
cd apps/api
swag init -g cmd/server/main.go -o docs/swagger
```

---

## Project Structure (Detailed)

See [`docs/architecture.md`](docs/architecture.md) for the full architecture reference and [`docs/adr/`](docs/adr/) for the decisions behind it. Related docs: [`docs/vision.md`](docs/vision.md), [`docs/requirements.md`](docs/requirements.md), [`docs/roadmap.md`](docs/roadmap.md), [`docs/database.md`](docs/database.md).

---

## Contributing

We welcome contributions from the community! Please read our [Contributing Guide](CONTRIBUTING.md) before submitting a pull request.

Key contribution areas:
- Bug fixes and feature requests via [GitHub Issues](https://github.com/foundryhq/foundryhq/issues)
- Frontend components and pages
- Backend API endpoints and business logic
- Documentation improvements
- Test coverage

---

## Roadmap

See [`docs/mvp.md`](docs/mvp.md) for what's shipping in Version 1 (Auth, Workspace/Team, Tasks) and [`docs/roadmap.md`](docs/roadmap.md) for the full sequencing after that — CRM, Meeting Notes, OKRs, KPI Dashboard, Notifications, then integrations (email/calendar sync, AI meeting summaries, Zapier/Make webhooks, a public API). White-label/self-hosted and SSO/SAML are explicitly iceboxed until a real customer asks.

---

## License

MIT © [FoundryHQ](https://github.com/foundryhq/foundryhq)
