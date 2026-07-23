# Contributing to FoundryHQ

Thank you for your interest in contributing to FoundryHQ! We're building the operating system for startups and we'd love your help. This guide will get you from zero to merged pull request.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Commit Messages](#commit-messages)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

---

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you agree to uphold this standard. Please report unacceptable behavior to the maintainers.

---

## How Can I Contribute?

### Good First Issues

Look for issues tagged [`good first issue`](https://github.com/foundryhq/foundryhq/labels/good%20first%20issue) — these are scoped, well-documented, and a great place to start.

### Ways to Contribute

- **Bug fixes** — Pick up any open bug from [Issues](https://github.com/foundryhq/foundryhq/issues)
- **Features** — Discuss in an issue first before building (to avoid wasted effort)
- **Documentation** — Improve the docs in the `docs/` directory or inline code comments
- **Tests** — Increase test coverage on the backend or frontend
- **Design** — UI/UX improvements (open an issue with mockups first)

---

## Development Setup

### 1. Fork and Clone

```bash
git clone https://github.com/YOUR_USERNAME/foundryhq.git
cd foundryhq
git remote add upstream https://github.com/foundryhq/foundryhq.git
```

### 2. Environment Variables

```bash
cp apps/api/.env.example apps/api/.env
# Edit .env with your local values
```

### 3. Start Infrastructure

```bash
docker compose up postgres -d
```

### 4. Backend

```bash
cd apps/api
go mod download
go run cmd/server/main.go
```

### 5. Frontend

```bash
cd apps/web
pnpm install
pnpm dev
```

### 6. Mobile (optional)

```bash
cd apps/mobile
pnpm install
npx expo start
```

---

## Project Structure

```
foundryhq/
├── apps/
│   ├── api/                    # Go backend
│   │   ├── cmd/server/         # main.go entrypoint
│   │   ├── internal/
│   │   │   ├── domain/         # Entity structs & repository interfaces
│   │   │   ├── usecases/       # Business logic (no HTTP, no DB details)
│   │   │   ├── repositories/   # DB implementations
│   │   │   └── handlers/       # Gin HTTP handlers
│   │   └── pkg/                # Shared utilities (jwt, db, logger)
│   ├── web/                    # React + TypeScript
│   └── mobile/                 # React Native + Expo
├── packages/
│   └── shared-types/           # Types shared between web and mobile
└── docs/                       # Documentation
```

**Backend architecture** follows Clean Architecture:
- `domain` → knows nothing about infrastructure
- `usecases` → depends only on `domain`
- `repositories` → implements domain interfaces
- `handlers` → calls usecases, speaks HTTP

---

## Coding Standards

### Go (Backend)

- Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `golangci-lint` before pushing: `cd apps/api && golangci-lint run`
- Format with `gofmt` / `goimports`
- Error wrapping: use `fmt.Errorf("context: %w", err)`
- No `panic` in handler code — always return errors

### TypeScript (Frontend / Mobile)

- Strict TypeScript — no `any`, no `@ts-ignore`
- React: functional components, hooks only
- State: Zustand for global, local state with `useState`
- Server state: React Query / TanStack Query
- Styling: Tailwind CSS utility classes; no inline styles
- Lint: `pnpm lint` must pass

### General

- No commented-out code
- No `TODO` comments without an associated GitHub issue number
- Meaningful variable names — code should read like prose

---

## Testing

### Backend

```bash
cd apps/api
go test ./...                    # All tests
go test ./internal/usecases/...  # Specific package
go test -cover ./...             # With coverage
```

- **Unit tests**: every usecase must have unit tests; mock repositories via interfaces
- **Integration tests**: located in `internal/repositories/postgres/*_test.go`; require a running DB

### Frontend

Testing isn't set up yet — neither `apps/web` nor `apps/mobile` has a test runner, no Vitest/React Testing Library dependency, and no `test` script. The only frontend checks today are lint and build:

```bash
pnpm --filter @foundryhq/web lint
pnpm --filter @foundryhq/web build

pnpm --filter @foundryhq/mobile lint
pnpm --filter @foundryhq/mobile build
```

This section will be updated with the runner, commands, and coverage expectations once frontend testing is added.

### What to Test

| Layer | Test Type | What to Cover |
|---|---|---|
| Usecases | Unit | Business logic, validation, error paths |
| Repositories | Integration | SQL correctness, constraint behavior |
| Handlers | Unit | Request binding, response shape, error codes |
| Components | Unit | Renders, user interactions, edge states |

---

## Pull Request Process

1. **Branch from `main`**
   ```bash
   git checkout main && git pull upstream main
   git checkout -b feat/your-feature-name
   ```

2. **Keep PRs focused** — one feature or fix per PR. Large PRs are hard to review.

3. **Update tests** — new behavior must have corresponding tests.

4. **Update docs** — if you added an endpoint, update the Swagger annotations.

5. **Pass CI** — all checks must be green before review.

6. **Fill in the PR template** — describe what changed and why.

7. **Request a review** — tag a maintainer or let the bot assign one.

8. **Address feedback** — push follow-up commits (don't force-push during review).

### PR Title Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(crm): add deal pipeline kanban view
fix(auth): refresh token not invalidated on logout
docs(api): document /meetings endpoints
test(tasks): add usecase unit tests for sprint creation
chore(deps): bump go to 1.22.4
```

---

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short summary>

[optional body]

[optional footer: Closes #123]
```

**Types:** `feat`, `fix`, `docs`, `test`, `chore`, `refactor`, `style`, `perf`

**Scopes:** `auth`, `crm`, `tasks`, `meetings`, `goals`, `kpi`, `team`, `api`, `web`, `mobile`, `infra`

---

## Reporting Bugs

Open a [GitHub Issue](https://github.com/foundryhq/foundryhq/issues/new?template=bug_report.md) with:

- **Environment** (OS, browser, API version)
- **Steps to reproduce**
- **Expected behavior**
- **Actual behavior**
- **Screenshots / logs** (if applicable)

---

## Suggesting Features

Open a [Feature Request Issue](https://github.com/foundryhq/foundryhq/issues/new?template=feature_request.md) with:

- **Problem statement** — what user pain does this solve?
- **Proposed solution**
- **Alternatives considered**
- **Mockups** (if UI-related)

Large features should be discussed in an issue before any implementation begins.

---

## Questions?

- Open a [Discussion](https://github.com/foundryhq/foundryhq/discussions)
- Join our community (link in README)

Thank you for making FoundryHQ better! 🚀
