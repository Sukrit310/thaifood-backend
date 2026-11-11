# ThaiFood Guide Backend Implementation Plan

**Goal:** Deliver every backend endpoint in API Contract v2, tested against PostgreSQL before a Render deployment using Supabase.

**Architecture:** Preserve the reference's auth/, handler/, store/, mw/ and main.go separation. PostgreSQL via pgxpool owns integrity; Go issues bcrypt-backed HS256 JWTs and reloads the member on protected requests so current database roles govern access. HTTP handlers depend on a typed store interface for fault tests; end-to-end tests exercise the real PostgreSQL store.

**Spec:** ThaiFood_Guide_Project_Brief_v2.md and ThaiFood_Guide_API_CONTRACT_v2.md. The supplied documents determine the design and authorize implementation. No source repository changes.

## Global Constraints

- net/http, pgxpool, JWT/bcrypt; environment DATABASE_URL, JWT_SECRET, ALLOWED_ORIGINS, PORT.
- Exact Contract v2 methods, paths, camelCase fields and nested error envelope.
- No reviews: rating null, reviewCount 0. Thai province/category name filters, page 1, limit 12, maximum 50.
- Register always member. Database current role governs admin actions. Owner/admin review deletion.
- Review integer 1–5 and at least 10 Unicode characters. Unique member/restaurant review and favorite.
- PostgreSQL and schema tests isolated from original MovieHub; never commit credentials.

## Implementation decisions

- Passwords 8–72 bytes (bcrypt maximum); display names 1–100 Unicode characters, normalized email.
- JWT one hour, explicit algorithm/issuer/audience/expiry checks. Missing/bad/expired tokens return 401.
- Admin PUT accepts partial restaurant fields (Contract says fields to edit); omitted fields preserved, imageUrl null clears image. Required create fields: name, provinceId, categoryId, address, description; coordinates optional but must be paired and in range.
- JSON unknown fields rejected, especially role in registration; one object per body, size limited to 1 MiB.
- Parameterized SQL, database uniqueness and cascading restaurant deletions; RLS enabled and public API database roles denied because all access is through Go.
- Seed data clearly labelled author-created fictional examples, 25 restaurants across 5 provinces, with source metadata; no scraped ratings/images.
- GET /health checks database connectivity, returns generic 503 on failure; unmatched paths and methods use the JSON error envelope.

## Tasks

- [x] 1. Store/schema: typed models/interface, all queries, idempotent schema/seed, database integration tests.
- [x] 2. Auth/middleware: password/JWT tests, strict tokens, member context, database-authorized roles, exact CORS allowlist, sanitized request logging.
- [x] 3. HTTP API: validation, all Contract v2 routes, standardized errors, pagination, owner/member/admin flows and fault tests.
- [x] 4. Runtime/tooling: bounded pool, timeouts, graceful shutdown, migrate/admin commands, local test PostgreSQL setup, Postman collection and CI.
- [x] 5. Real PostgreSQL API integration: guest/auth/profile/password/search/review/favorite/admin paths, errors, injection, stale role, constraints and cascade.
- [x] 6. Review and verification: gofmt, go vet, go test -race, build, Postman runner and secret review.
- [x] 7. Deployment: confirm target Supabase and Render workspace, migrate/seed, test Supabase-backed API, push feature branch, create/update Render service, wait live and smoke test. If credentials/target remain unavailable, report the specific external prerequisite without claiming deployment.

## Review Focus

- Invalid/expired/forged tokens and a previously admin user demoted in the DB.
- Unicode Thai review length, malformed/oversized/trailing JSON, invalid pagination and coordinates.
- Concurrent duplicate reviews/favorites/emails; empty arrays must encode as [] rather than null.
- SQL metacharacters in search and case-insensitive duplicate emails.
- Role isolation, foreign-key validation, aggregate recomputation after deletes, cascading cleanup, DB failures without leaking credentials.

## Execution record

- Tasks1–6 complete: PostgreSQL17 real API/store tests,71 test/subtest passes with race detector,64 Newman requests/296 assertions, zero failures, vet/build/Docker/schema validation passed.
- Independent reviews: invalid UTF-8 filter finding fixed with failing-then-passing regression; runtime pgx mode verified by Newman against built binary.
- Security scan: Go1.27.3 standard-library finding fixed by toolchain/Docker/CI Go1.27.4; final govulncheck clean.
- Task 7 complete (2026-11-12 Asia/Bangkok): the user supplied the Supabase project/pooler credentials and explicitly authorized transfer to the confirmed Render service. Migration/seed and 8 integration tests passed on Supabase; predeploy Postman passed 64 requests / 296 assertions. Render deploy dep-dbqacjijnfac73dbf24g of commit 53cdab8 is live at https://thaifood-guide-api-a7c7.onrender.com. Live Postman passed 64 requests / 294 assertions with one explicitly identified upstream WAF 403 SQL probe; the default strict Go/Supabase probe still requires 200 JSON and passes. Removed 7 recorded temporary accounts and verified final health 200. The original MovieHub repo/service is unchanged. The user deferred the Vercel origin; localhost CORS is configured and the handoff documents the next step.
