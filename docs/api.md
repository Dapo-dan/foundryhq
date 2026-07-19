# API Reference

Human-curated overview of FoundryHQ's REST API. This is the index — the generated, exhaustive spec lives at `swagger/` (produced by `swag init`, served at `/swagger/index.html`) per the annotation standard in `../.ai/documentation/api.md`. If this file and the generated spec ever disagree, the generated spec is correct — file an issue against this doc.

## Base URL & Versioning

- **Dev:** `http://localhost:8080`
- **Auth:** `Authorization: Bearer <access_token>` on every route except `auth/*` and the health check
- Breaking changes are versioned with a path prefix (`/v2/...`); additive changes ship without a version bump — see `../.ai/documentation/api.md#versioning`

## Response Envelope

```json
// Success
{ "data": { ... } }

// Success (paginated list)
{ "data": [ ... ], "meta": { "page": 1, "per_page": 20, "total": 143 } }

// Error
{ "error": { "code": "validation_error", "message": "title is required", "field": "title" } }
```

## Endpoint Index

### Operational
| Method | Path | Description |
|---|---|---|
| GET | `/health` | Liveness check — process is up, no dependency checks |
| GET | `/ready` | Readiness check — additionally verifies the database is reachable |
| GET | `/version` | Running build's version and commit |

### Auth
| Method | Path | Description |
|---|---|---|
| POST | `/auth/register` | Email/password registration |
| POST | `/auth/login` | Email/password login |
| POST | `/auth/oauth/{provider}` | OAuth login (`google`, `github`) |
| POST | `/auth/refresh` | Exchange refresh token for a new access token |
| POST | `/auth/logout` | Invalidate the current refresh token |

### Workspaces & Team
| Method | Path | Description |
|---|---|---|
| GET | `/workspaces/{id}` | Get workspace details |
| PATCH | `/workspaces/{id}` | Update workspace name/logo/slug |
| GET | `/workspaces/{id}/members` | List members and roles |
| POST | `/workspaces/{id}/members/invite` | Invite by email |
| PATCH | `/workspaces/{id}/members/{memberId}` | Change a member's role |

### CRM
| Method | Path | Description |
|---|---|---|
| GET / POST | `/contacts` | List / create contacts |
| GET / PATCH / DELETE | `/contacts/{id}` | Read / update / soft-delete a contact |
| GET / POST | `/companies` | List / create companies |
| GET / POST | `/deals` | List (Kanban-grouped by stage) / create deals |
| GET / PATCH / DELETE | `/deals/{id}` | Read / update (including stage change) / soft-delete a deal |
| GET / POST | `/deals/{id}/activities` | List / log an activity against a deal |

### Tasks & Sprints
| Method | Path | Description |
|---|---|---|
| GET / POST | `/tasks` | List (backlog or filtered) / create tasks |
| GET / PATCH / DELETE | `/tasks/{id}` | Read / update (status, assignee, priority) / soft-delete a task |
| GET / POST | `/sprints` | List / create sprints |
| GET | `/sprints/{id}` | Sprint detail with tasks grouped by status |
| GET | `/sprints/{id}/velocity` | Computed velocity for the sprint |

### Meetings
| Method | Path | Description |
|---|---|---|
| GET / POST | `/meetings` | List / create meetings |
| GET / PATCH | `/meetings/{id}` | Read / update notes and attendees |
| POST | `/meetings/{id}/action-items` | Create an action item from highlighted text |

### Goals (OKRs)
| Method | Path | Description |
|---|---|---|
| GET / POST | `/objectives` | List (alignment tree) / create objectives |
| GET / PATCH | `/objectives/{id}` | Read / update an objective |
| POST | `/objectives/{id}/key-results` | Add a key result |
| POST | `/key-results/{id}/check-in` | Record a progress update |

### KPIs
| Method | Path | Description |
|---|---|---|
| GET / POST | `/kpis` | List / define KPIs |
| POST | `/kpis/{id}/values` | Record a time-series value |
| GET | `/kpis/{id}/values` | Time-series history |
| POST | `/kpis/snapshot` | Generate a read-only shareable snapshot link |

### Notifications
| Method | Path | Description |
|---|---|---|
| GET | `/notifications` | List for the current user |
| PATCH | `/notifications/{id}/read` | Mark as read |
| PATCH | `/notifications/preferences` | Update digest preferences |

## Error Codes

| Code | HTTP status | Meaning |
|---|---|---|
| `validation_error` | 400 | Request body failed validation; `field` is set when applicable |
| `unauthorized` | 401 | Missing or invalid access token |
| `forbidden` | 403 | Authenticated, but the workspace role doesn't permit this action |
| `not_found` | 404 | Resource doesn't exist, or exists in a different workspace |
| `conflict` | 409 | e.g., inviting an email already a member |
| `rate_limited` | 429 | Too many requests from this client (e.g. `/auth/login` attempts) — retry after a delay |
| `internal_error` | 500 | Unexpected server error — logged with a request ID for correlation |

## Adding an Endpoint

Follow the annotation and envelope standard in `../.ai/documentation/api.md` — every new handler needs complete `godoc` swaggo annotations, and this index should get a one-line entry in the same PR.
