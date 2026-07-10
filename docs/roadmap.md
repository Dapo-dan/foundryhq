# Roadmap

*Last updated: 2026-07-10.* This is the current snapshot of what FoundryHQ is building and why. For the prioritization framework and review cadence behind these decisions, see `../.ai/startup/roadmap.md`. For the reasoning behind the product bet itself, see `vision.md`.

## Shipped

The current GA feature set — see `../README.md#features` for the user-facing description of each:

- Authentication (email/password, OAuth via Google/GitHub, JWT + refresh tokens, role-based access)
- CRM (contacts, companies, deal pipeline, activity timeline)
- Task Management (Kanban/list views, sprints, sub-tasks, labels, priorities)
- Meeting Notes (structured notes, action item extraction, linking to contacts/projects)
- Goal Tracking / OKRs (company/team/personal objectives, key result check-ins, alignment tree)
- Team Management (invites, roles, activity feed)
- KPI Dashboard (custom KPIs, time-series values, goal alignment)
- Notifications (in-app, push, email, @mentions)

## Now (current sprint)

- Hardening the CRM activity timeline for performance at higher deal volume (Early Sales Rep pain point: losing context on cold deals)
- Investor/Advisor read-only snapshot sharing — shareable link, no full workspace login required

## Next (1–2 sprints out)

- **Email integration (IMAP/SMTP sync)** — logs emails to the CRM activity timeline automatically; highest-leverage fix for the Early Sales Rep persona's "reduce clicks to log an activity" need
- **Calendar integration (Google Calendar, Outlook)** — auto-creates meeting records from calendar events, reducing meeting-notes setup friction

## Later (3+ sprints)

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
