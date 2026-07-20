# Roadmap

*Last updated: 2026-07-10.* This is the current snapshot of what FoundryHQ is building and why. For the prioritization framework and review cadence behind these decisions, see `../.ai/startup/roadmap.md`. For the reasoning behind the product bet itself, see `vision.md`.

## Shipped

Nothing yet — `apps/` and `packages/` are empty scaffolding, no application code has been written. The feature descriptions in `../README.md#features` and `../PROJECT.md` describe the target state, not current state.

## Now (Version 1 — this month)

The MVP scope — see `mvp.md` for the full checklist and definition of done:

- Authentication (email/password, JWT + refresh tokens) — REQ-01
- Workspace & Team (create workspace, invite members, switch workspace, Owner/Member roles) — REQ-02 (partial)
- Task Management (Kanban board, CRUD tasks, assign users, status updates) — REQ-11
- Mobile (React Native/Expo, scaffolded in `apps/mobile`) — ships Auth, Workspace/Team, and Tasks alongside web, not after; full parity for later modules deferred per REQ-NFR-05

## Next (v1.1 — after MVP ships)

- **CRM** (contacts, companies, deal pipeline, activity timeline) — REQ-03, REQ-04
- **OAuth** (Google, GitHub) and full role set (Admin, Viewer) — completes REQ-01, REQ-02
- **Task Management (advanced)** — sprints, velocity tracking, sub-tasks, labels, priorities — REQ-05
- **Meeting Notes** (structured notes, action item extraction, linking to contacts/projects) — REQ-06

## Later (v1.2+)

- **Goal Tracking / OKRs** (company/team/personal objectives, key result check-ins, alignment tree) — REQ-07
- **KPI Dashboard** (custom KPIs, time-series values, goal alignment) — REQ-08
- **Investor/Advisor read-only snapshot sharing** — REQ-09
- **Notifications** (in-app, push, email, @mentions) — REQ-10
- **Email integration (IMAP/SMTP sync)** — logs emails to the CRM activity timeline automatically; highest-leverage fix for the Early Sales Rep persona's "reduce clicks to log an activity" need
- **Calendar integration (Google Calendar, Outlook)** — auto-creates meeting records from calendar events, reducing meeting-notes setup friction
- **AI meeting summaries** — auto-drafted summaries and suggested action items from meeting notes
- **Zapier / Make webhooks** — lightweight integration surface for teams with tools outside FoundryHQ's core set
- **Public API for integrations** — a scoped, documented external API (distinct from the internal API in `api.md`), gated on real demand from teams already embedded in the product

## Icebox

Explicitly considered and deferred — see the anti-patterns in `../.ai/startup/roadmap.md` before re-litigating these:

- **White-label / self-hosted enterprise edition** — premature enterprise feature; revisit only when a real enterprise prospect asks, not speculatively
- **SSO / SAML** — same reasoning; no current customer has asked
- **Custom fields / workflow builder** — adds configuration surface area that cuts against the "radical simplicity" principle in `vision.md` until a specific, validated need emerges

## How Items Move

An item only moves from **Next**/**Later** into **Now** after it's been scoped against `requirements.md` (REQ-IDs, acceptance criteria) and reviewed for architectural impact against `architecture.md`. An item only moves into **Shipped** once it's live in production and its `../README.md` feature entry is updated in the same change — see `../.ai/documentation/readme.md`.
