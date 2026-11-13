# Admin member roles — Backend implementation plan

> **For agentic workers:** Use superpowers:executing-plans to implement these tasks in this workspace.

**Goal:** Implement only the backend requirements in `ThaiFood_Admin_Role_Management_Codex.md`, sections 3 and 7.

**Architecture:** Extend the existing Go handler, admin middleware and repository. A short READ COMMITTED transaction takes a SHARE ROW EXCLUSIVE lock on members, reloads the actor's role, checks self/last-admin restrictions, and updates with SQL parameters. This serializes role changes across API instances and keeps ordinary reads available; member writes wait briefly. No schema or authentication redesign is required.

**Tech stack:** Existing Go net/http, pgx and PostgreSQL; no new dependencies.

**Spec:** `ThaiFood_Admin_Role_Management_Codex.md` (backend scope only).

## Global constraints

- PATCH /api/admin/members/{id}/role uses the existing admin middleware.
- Only the exact JSON key `role`, with `member` or `admin`, is accepted.
- Return the existing four-field Member object; no hash, token or secret.
- Self changes return 403 FORBIDDEN, last-admin protection returns 409 LAST_ADMIN; identical roles return 200 without rewriting updated_at.
- Reload actor authorization inside the transaction to reject requests admitted before concurrent demotion.
- Work in the user's current repository as requested. Leave the supplied untracked specification untouched. No production DB, frontend, push or deployment changes.
- Test only against a disposable local database via TEST_DATABASE_URL.

## Review focus

- Two admins concurrently demoting each other must leave exactly one admin.
- A request waiting behind a committed demotion must lose admin authority.
- Invalid bodies (including null, extra keys and differently cased keys) must not mutate data.
- Same-role requests must remain subject to authorization and self-change rules.
- Existing JWTs must gain/lose admin access on the next request without login.

## Tasks

### 1. Endpoint and validation

Files: `handler/handler.go`, `handler/auth.go`, `handler/member_roles_test.go`, `store/models.go`.

Interface: `UpdateMemberRole(context.Context, actorID int64, targetID int64, role string) (Member, error)`.

- [x] Write table-driven HTTP tests for authorization, success, exact input validation, self restrictions, response shape and sanitized store errors.
- [x] Run tests and observe missing endpoint failures.
- [x] Register the admin route, implement strict role parsing using the existing body helper, and call the repository with the authenticated actor ID.
- [x] Confirm focused handler tests pass.

### 2. Transaction and database-backed coverage

Files: `store/members.go`, `store/member_roles_test.go`, `integration/member_roles_test.go`.

- [x] Add real database tests for persistence, idempotence, forbidden callers, self/last-admin safety, rollback and concurrent cross-demotions.
- [x] Add HTTP integration coverage for promotion/demotion with existing tokens, direct member/guest rejection, invalid input, registration/profile escalation rejection, and concurrent cross-demotions.
- [x] Observe tests failing before implementing the store method.
- [x] Implement the transaction, current actor check, last-admin count and parameterized update; use existing error mapping plus LAST_ADMIN.
- [x] Run focused store/integration tests against disposable PostgreSQL.

### 3. Contract, handoff and verification

Files: `API_CONTRACT.md`, `README.md`, `docs/ADMIN_ROLE_MANAGEMENT.md`, Postman collection/README.

- [x] Document request/response/errors, idempotence, lock tradeoff and FE test instructions.
- [x] Add a Postman sequence operating only on collection-generated members and restore their roles.
- [x] Run `go test -race -count=1 ./...` (with `-json` for counting), `go vet ./...`, `go build ./...`, gofmt and diff checks.
- [x] Review final changes against backend acceptance criteria and report actual pass/fail/skip counts and deployment limitations.

## Execution record

- Baseline suite passed; DB tests initially skipped before the disposable database was set up.
- RED: handler and real HTTP tests returned 404 for the missing endpoint; store tests could not compile because the new method did not exist.
- GREEN: focused handler/store/integration tests passed after implementation.
- Full verification: 39 top-level tests + 90 subtests passed with race detection; zero failed/skipped tests. Vet/build/format/diff checks passed.
- Independent review found no material code/test findings. Documentation/Postman review caught `adminId` capture in a prerequest script; a lifecycle execution reproduced the TypeError, then passed after moving capture to the response test script.
- Postman collection has 78 requests and 82 syntactically valid scripts. Full Newman and browser/production validation are not claimed.
- Work remains uncommitted in the requested workspace. No frontend, schema, production or deployment changes.
