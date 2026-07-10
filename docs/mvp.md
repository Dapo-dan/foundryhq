# MVP — Version 1

*Last updated: 2026-07-10.* FoundryHQ's full vision (see `vision.md`) spans eight modules — too much for a single month. This document scopes exactly what ships in **Version 1** versus what's deferred. Nothing in this repo is built yet (see `../PROJECT.md` and `../README.md` for target-state architecture); this is the build order.

## Why this scope

The smallest usable loop is: a user can sign up, create a workspace, bring in their team, and manage shared work. That's Auth + Workspace/Team + Tasks. Everything else (CRM, Meeting Notes, OKRs, KPI Dashboard, Notifications) depends on that loop existing first and is deferred to v1.1+.

## Version 1

### Authentication
- [ ] Register (email + password)
- [ ] Login
- [ ] Logout
- [ ] JWT authentication (access + refresh tokens)

*Deferred to v1.1+: OAuth (Google, GitHub) — see REQ-01.*

### Workspace & Team
- [ ] Create workspace
- [ ] Invite members (via email)
- [ ] Switch workspace
- [ ] Basic roles: Owner, Member

*Deferred to v1.1+: Admin/Viewer roles, read-only access, shareable snapshots — see REQ-02.*

### Tasks
- [ ] Kanban board (view by status)
- [ ] CRUD tasks (create, edit, delete)
- [ ] Assign users
- [ ] Status updates

*Deferred to v1.1+: sprints, velocity tracking, sub-tasks, labels, priorities, due dates — see REQ-05, REQ-11.*

## Version 1.1+ (everything else)

Moved out of v1 scope entirely, in rough sequence — see `roadmap.md` for the live prioritization:

- **CRM** — contacts, companies, deal pipeline, activity timeline (REQ-03, REQ-04)
- **Meeting Notes** — structured notes, action item extraction (REQ-06)
- **Goal Tracking (OKRs)** — objectives, key results, alignment tree (REQ-07)
- **KPI Dashboard** — custom KPIs, time-series tracking, shareable snapshots (REQ-08, REQ-09)
- **Notifications** — in-app, push, email, digests (REQ-10)
- **Team Management (advanced)** — Admin/Viewer roles, workspace settings (part of REQ-02)
- **Task Management (advanced)** — sprints, velocity, sub-tasks, labels (REQ-05)

## Out of scope for v1 (explicitly)

- Any non-functional requirement beyond what's needed to ship v1 safely (see `requirements.md` for the full NFR list — REQ-NFR-01 and REQ-NFR-03 apply to v1; REQ-NFR-02, 04, 05 apply once the relevant modules ship)
- Mobile app — v1 is web-only
- Everything in `roadmap.md`'s "Icebox" (white-label/self-hosted, SSO/SAML, custom fields/workflow builder)

## Definition of done for v1

- [ ] A new user can register, verify they're logged in, and log out
- [ ] A logged-in user can create a workspace and invite a second user by email
- [ ] Both users can see, create, edit, delete, and reassign tasks on a shared Kanban board scoped to their workspace
- [ ] REQ-01 (Auth) and REQ-11 (basic task management) are marked "Shipped" in `requirements.md`'s traceability table, and this checklist is fully checked
