# Vision

## The Problem

Most startups run on 10+ disconnected tools — a CRM here, a task tracker there, a notes app somewhere else, a spreadsheet holding it all together. Context gets lost between tools, founders and small teams waste time on manual busywork (copying meeting notes into a task tracker, re-explaining deal context because the CRM doesn't talk to Slack), and as the team grows, nobody has a single, trustworthy view of how the company is actually doing.

This isn't a tooling-quality problem — Notion, Linear, and HubSpot are all good at what they do individually. It's an integration problem: nothing in a founding team's stack was built assuming the same five people are running sales, product, and operations simultaneously.

## The Vision

**FoundryHQ is the unified operating system built specifically for startups.** One workspace where a team tracks customers (CRM), ships product (tasks and sprints), runs meetings (notes and action items), and measures what matters (OKRs and KPIs) — from day zero to Series B.

The goal is radical simplicity: everything a founding team needs, nothing they don't, in a single fast workspace that grows with them. FoundryHQ is not trying to out-feature Salesforce or Jira — it's trying to remove the *seams* between the tools a 2–15 person team already reaches for.

## Design Principles

1. **First value in under 3 minutes.** No forced configuration before a new workspace can touch real data (see the onboarding flow in `../.ai/product/user-flows.md`).
2. **Connected by default, not by integration.** A meeting action item becomes a task automatically. A deal's activity log is the same activity log the founder sees on the company dashboard. The value isn't any single feature — it's that they share one data model.
3. **Speed over completeness.** Founders use FoundryHQ between back-to-back calls. Every core interaction (log an activity, update a task, check a KPI) must be fast, not merely possible.
4. **Progressive disclosure.** Advanced configuration (custom fields, complex automations) stays out of view until a team is big enough to need it.
5. **Prunable, not permanent.** Features are validated against real usage before they're allowed to add complexity — see the anti-patterns in `roadmap.md`.

## Who This Is For

FoundryHQ is built around four core personas — see `../.ai/product/user-personas.md` for the full detail behind each:

| Persona | What they need most |
|---|---|
| **Founder-Operator** | One place to see what's happening across the company, without a sales or ops background |
| **Product Engineer** | Sprint work tied to real goals, without PM ceremony overhead |
| **Early Sales Rep** | A pipeline lighter than Salesforce, sturdier than a spreadsheet |
| **Investor/Advisor** | A trustworthy, read-only signal on portfolio company health |

## Non-Goals

- FoundryHQ is not an enterprise suite. SSO, audit logs, and SOC2 are deliberately deferred until a real enterprise customer asks for them (see `roadmap.md` anti-patterns).
- FoundryHQ does not aim to replace best-in-class point tools for large teams (e.g., a 200-person sales org's need for a heavyweight CRM). The bet is the founding-team stage, not "CRM for everyone."
- FoundryHQ does not chase competitor feature parity as a strategy — see the roadmap philosophy in `roadmap.md`.

## Where This Fits

- **What we're building right now:** `roadmap.md`
- **Why the system is built the way it is:** `architecture.md` and `adr/`
- **What "done" means for a given feature:** `requirements.md`
