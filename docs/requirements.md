# Requirements

Current, high-level requirements for FoundryHQ's GA feature set, prioritized with MoSCoW. This is the traceability anchor between `vision.md`/`roadmap.md` (why) and individual tickets/PRs (what got built). For the elicitation process and template behind these entries, see `../.ai/business-analysis/requirements.md`.

## REQ-01: Workspace-scoped authentication

**Type:** Functional
**Priority:** Must
**Persona:** All

The system shall authenticate users via email/password or OAuth (Google, GitHub) and scope every subsequent request to the authenticated user's workspace membership.

**Rationale:** Every other feature depends on knowing who the user is and which workspace's data they can see — this is the load-bearing requirement for the whole product.

## REQ-02: Role-based access control

**Type:** Functional
**Priority:** Must
**Persona:** Founder-Operator (grants access), all personas (subject to access)

The system shall support four workspace roles — Owner, Admin, Member, Viewer — with Viewer granted read-only access across all modules.

**Rationale:** Directly enables the Investor/Advisor persona's read-only snapshot need without requiring a second product surface.

## REQ-03: Deal pipeline with Kanban view

**Type:** Functional
**Priority:** Must
**Persona:** Early Sales Rep

The system shall let a user create a deal with a name, linked contact/company, stage, and value, and display deals grouped by stage in a Kanban view.

**Rationale:** Core CRM value proposition — replacing a spreadsheet with something that shows pipeline state at a glance.

## REQ-04: Activity timeline per deal/contact

**Type:** Functional
**Priority:** Must
**Persona:** Early Sales Rep

The system shall record calls, emails, and notes against a deal or contact and display them in chronological order on the relevant detail page.

**Rationale:** Directly answers "what happened last on this account" — the Early Sales Rep's stated core pain in `../.ai/product/user-personas.md`.

## REQ-05: Sprint planning with velocity tracking

**Type:** Functional
**Priority:** Must
**Persona:** Product Engineer

The system shall let a user move backlog tasks into a time-boxed sprint and calculate velocity as the sum of story points completed within the sprint's date range.

**Rationale:** Connects day-to-day engineering work to a measurable, reviewable outcome without requiring a separate PM tool.

## REQ-06: Meeting notes with action item extraction

**Type:** Functional
**Priority:** Must
**Persona:** All (primarily Founder-Operator, Product Engineer)

The system shall let a user highlight text within a meeting note and convert it into an action item, automatically assigned to the selected team member's task list.

**Rationale:** Closes the specific gap called out in `vision.md` — meeting context and follow-up work living in two disconnected tools.

## REQ-07: OKR alignment tree

**Type:** Functional
**Priority:** Must
**Persona:** Founder-Operator, Product Engineer

The system shall let teams and individuals align their objectives to a parent (team or company) objective, and roll up progress in an alignment tree view.

**Rationale:** Gives the Founder-Operator persona the "one place to see what's happening" view without manually cross-referencing separate goal docs per team.

## REQ-08: KPI dashboard with time-series tracking

**Type:** Functional
**Priority:** Must
**Persona:** Founder-Operator, Investor/Advisor

The system shall let a user define a KPI with a target (number, percent, or currency) and record values over time, visualized against the target.

**Rationale:** The primary artifact behind the Investor/Advisor's shareable snapshot (REQ-09) and the Founder-Operator's company-health view.

## REQ-09: Read-only shareable snapshot

**Type:** Functional
**Priority:** Should
**Persona:** Investor/Advisor

The system shall let a workspace Owner/Admin generate a shareable, read-only link to a KPI/OKR snapshot that does not require the recipient to hold a full workspace login.

**Rationale:** The Investor/Advisor persona explicitly does not want a full-login workspace experience — see `../.ai/product/user-personas.md`.

## REQ-10: Notification delivery across channels

**Type:** Functional
**Priority:** Should
**Persona:** All

The system shall deliver notifications in-app, via push (mobile), and via email, with per-user digest preferences and @mention support across tasks, meetings, and goals.

**Rationale:** Supports asynchronous, distributed founding teams who aren't always in the app at the same time.

## Non-Functional Requirements

| ID | Requirement | Priority |
|---|---|---|
| REQ-NFR-01 | API responses for primary list views (deals, tasks, KPIs) return in under 300ms at p95 under normal load | Must |
| REQ-NFR-02 | Web UI meets WCAG 2.1 AA — see `../.ai/design/accessibility.md` | Must |
| REQ-NFR-03 | Access tokens expire in 15 minutes; refresh tokens in 7 days | Must |
| REQ-NFR-04 | Mutations apply optimistically in the UI and roll back cleanly on server error | Should |
| REQ-NFR-05 | Mobile parity for CRM, Tasks, and Notifications; OKR/KPI dashboards may remain web-only in the interim | Could |

## Traceability

| REQ-ID | Feature area | Persona | Status |
|---|---|---|---|
| REQ-01 | Authentication | All | Shipped |
| REQ-02 | Authentication / Team Management | All | Shipped |
| REQ-03 | CRM | Early Sales Rep | Shipped |
| REQ-04 | CRM | Early Sales Rep | Shipped |
| REQ-05 | Task Management | Product Engineer | Shipped |
| REQ-06 | Meeting Notes | Founder-Operator, Product Engineer | Shipped |
| REQ-07 | Goal Tracking (OKRs) | Founder-Operator, Product Engineer | Shipped |
| REQ-08 | KPI Dashboard | Founder-Operator, Investor/Advisor | Shipped |
| REQ-09 | KPI Dashboard | Investor/Advisor | In progress — see `roadmap.md` "Now" |
| REQ-10 | Notifications | All | Shipped |

New requirements should be added here with a fresh REQ-ID before a PRD or ticket is written — see the elicitation workflow in `../.ai/business-analysis/requirements.md` and the PRD template in `../.ai/product/prd.md`.
