# Kheun Lift — Go backend

Powerlifting leaderboard app. Users log Squat / Bench / Deadlift lifts with a
video link. Lifts are ranked by exercise and weight class, and only `verified`
lifts appear on the leaderboard.

This repo is the **new backend that will be deployed** (Go, Gin, pgx, goose).
Sibling folders under `../` (each its own git repo):

| Folder | Role |
|---|---|
| `../kheun-lift-frontend-next` | New frontend, will be deployed (Next.js 16, still uses demo data) |
| `../kheun-lift-backend` | Old bootcamp backend (Express, Prisma), **reference only** |
| `../Kheun-Lift-frontend` | Old bootcamp frontend (React, Vite), **reference only** |

- `PLAN.md` here is the main plan (decisions, architecture, backend tasks by
  Phase 1–7, migrations, env). `../kheun-lift-frontend-next/PLAN.md` holds
  frontend tasks with the same Phase numbers. Tick tasks as they are done.
- Do not edit the two old repos unless the user asks. Use them to see how
  features worked before.

## Rules

### 1. Require approval before making changes
Before modifying code or adding anything to the project, show the user
the proposed changes, affected files, and reasons.
Wait for explicit approval before creating, editing, deleting, or
renaming files. If the scope changes, request approval again.

Read-only inspection is allowed before approval.

### 2. Require approval for database changes
Never modify database schemas or data without explicit approval.
This includes running migrations, seed scripts, and any commands
that write to the database.

Explain the intended operation and its impact before requesting approval.
Approval for code changes does not imply approval for database changes.

### 3. Tests only for risky logic, with approval
Create test files only when the user approves them as part of the change.
Keep tests small and focused on risky logic: auth, passwords, OTP,
weight class lookup, leaderboard ranking. No tests for simple code.

### 4. Review with Ponytail before making changes
Always use the Ponytail plugin's applicable skills to review relevant
code before proposing changes.

Identify unnecessary code, duplication, and over-engineering.
Propose the smallest change that preserves existing logic and behavior.
Include any proposed cleanup in the approval request; do not remove
code without approval.

If Ponytail is unavailable, inform the user before proceeding.

### 5. Do not run builds automatically
Do not run build commands unless the user explicitly requests a build.
This includes builds for validation, reviews, or completing a task.

### 6. Git is the user's job
Never commit, stage, push or move files with git unless the user asks in
that message. When work is ready, suggest a short commit message and list
its files. If the work should be split, give one message and file list
per commit, in order. The user commits.

## How the apps fit together

- This API runs on port `3001`; Next runs on `3002` (port 3000 is taken on
  this machine).
- Auth uses an httpOnly cookie, `kheun_session` (SameSite=Lax). There is no JWT.
- Next calls this API directly, cross-origin, at `NEXT_PUBLIC_API_URL`.
  Prod: Next on `app.<domain>`, Go on `api.<domain>` (same site).
- Exactly one origin is allowed, `APP_ORIGIN` (e.g. `http://localhost:3002`):
  CORS middleware (`internal/middleware/cors.go`, with credentials) plus
  `AddTrustedOrigin` for writes. Any other origin still gets 403 on
  POST/PATCH/DELETE. Locally set `COOKIE_SECURE=false`.
- Cross-subdomain cookies on prod need the cookie `Domain` (Phase 1, not done yet).
- Endpoints match the old Express API: `/auth`, `/users`, `/exercises`,
  `/lifts`, `/leaderboards`, `/weight-classes`. One difference: the
  `/leaderboards/exercises/...` list returns `{message, data}`, not a bare array.
- Mock users and lifts live in `seeds/` and run only via `make seed` (local
  dev). `make migrate` never loads them.
- No endpoint exists yet to verify or reject lifts.
