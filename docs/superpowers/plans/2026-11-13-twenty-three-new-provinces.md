# Twenty-three new provinces implementation plan

> **For agentic workers:** Use superpowers:executing-plans for implementation. Independent source research is delegated under dispatching-parallel-agents; the main worker owns SQL/import/verification.

**Goal:** Add 23 real businesses, one in each of 23 previously empty provinces, with matching credited photographs, bringing 30 restaurants / 19 provinces to 53 / 42.

**Architecture:** Preserve the existing API/schema and all thirty restaurant records. Append a source-deduplicated INSERT block to the seed, then apply only that block to the live ThaiFood database after backup and transactional baseline checks.

**Tech stack:** Existing Go/pgx/PostgreSQL API; public TAT/government/operator/Commons facts and individually licensed Commons photographs.

**Spec:** User request on 2026-11-13: เพิ่มข้อมูลร้านจริงอีก 23 ร้าน ในจังหวัดที่ไม่มีข้อมูล. Follow the same factual-data and photo-credit requirements in `docs/RESTAURANT_DATA_SOURCES.md`.

## Constraints and review focus

- No invented address components, venue coordinates, reviews or current opening status. Match the exact business/branch and province to its photograph.
- All 23 target provinces must be distinct and absent from the saved live baseline.
- Keep original rows and user data unchanged; do not execute earlier retirement blocks during this additive import.
- Preserve the existing `.env.example` modification; exclude credentials, raw datasets and private backup artifacts from commits.
- The total exceeds the API's maximum page size of 50. Verify all 53 results across pages, including photo credits for rows after the first fifty.

## Execution

- [x] Research 23 distinct province/business/photo matches; retain factual and photo metadata evidence locally.
- [x] Update seed replay/store/API expectations to 53 restaurants and 42 covered provinces. Verify the old seed fails the new count check.
- [x] Append only the 23 reviewed records to `seed.sql`. Update README, FE handoff and source register; document missing locality details honestly.
- [x] Run formatting, vet, all race tests with local PostgreSQL 17, build, and independent final review.
- [x] Back up the thirty restaurant rows with mode 0600. Transactionally apply the INSERT-only block and replay; require unchanged old fields and member/review/favorite counts.
- [x] Verify live list/detail/province filters, default and maximum-size pagination, reference data, health, all 23 image URLs and existing thirty API rows.
- [x] Record actual results in `docs/TWENTY_THREE_NEW_PROVINCES_IMPORT_RESULTS.md`. The verified change is ready for commit/push; subsequent CI status is tracked on the feature branch. No new merge/deployment was requested.
