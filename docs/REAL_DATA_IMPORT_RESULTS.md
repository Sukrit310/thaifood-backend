# Live real-restaurant import — 12 November 2026

The existing Render API now serves **20 real restaurants in 9 provinces**, each with a different, matched photograph and visible attribution. These are 19 factual records from the TAT open dataset and one from Breezeo's official operator. See [the full source and image-license register](RESTAURANT_DATA_SOURCES.md).

## Production database result

Target: the existing ThaiFood Supabase project `xorsagxmphehnjurquen`, used by `https://thaifood-guide-api-a7c7.onrender.com`. The import used the already configured TLS database connection without exposing credentials.

| Check | Before | After first import | After replay |
| --- | ---: | ---: | ---: |
| Restaurants | 25 | 20 | 20 |
| Explicit fictional seed rows | 25 | 0 | 0 |
| Restaurants with an image URL | 0 | 20 | 20 |
| Members | 2 | 2 | 2 |
| Reviews | 0 | 0 | 0 |
| Favorites | 0 | 0 | 0 |

A backup of the 25 original fictional restaurant rows was saved locally with file permission `0600` before writing; it contains no member/review/favorite contents and is not committed. A fresh preflight confirmed none of those rows had been edited or referenced. Only those unused, untouched fictional records were retired. Both seed runs completed successfully. Restaurant IDs and the fingerprint of every restaurant field were identical after the second run, confirming replay did not duplicate or rewrite records.

The API application code was not changed or redeployed. It queries PostgreSQL for each request, so the database change is already visible at the existing public URL. No Moviehub, FE repository, CORS, member account or authorization changes were made.

## Verification

- `go vet ./...`: passed.
- `go test -race -count=1 -json ./...`: **73 tests/subtests passed, 0 failed, 0 skipped**, including real PostgreSQL store/API integration in isolated schemas that were removed afterward.
- `go build -trimpath`: passed.
- `gofmt` and `git diff --check`: passed.
- The new real-data regression first failed against the old seed because its first restaurant had no photograph; it passes with the real seed.
- The preservation regression also reproduced deletion of an admin-edited old sample before the timestamp guard was added. It now verifies retention of reviews, favorites, admin-created rows, edits to real records and admin-edited old samples, while an unused untouched sample is removed.
- The direct-write coordinate test now explicitly writes `longitude=NULL` when testing an unpaired latitude. The old fixture assumed every seed coordinate was null; real seed records can have coordinates. The database constraint itself is unchanged.
- A separate read-only reviewer checked the SQL, source mapping, photo/license associations and preservation behavior; no outstanding findings remained.

Live public API verification completed on **12 November 2026 at approximately 22:21 Asia/Bangkok**:

- `GET /api/restaurants?limit=50`: 20 items, all with a different `imageUrl`, real names and unchanged Contract v2 JSON keys.
- All 20 detail endpoints matched the list and the reviewed records. Floating-point coordinates were compared with a small numeric tolerance; the largest serialization/conversion difference was approximately `2.85e-14` degrees.
- Default pagination returned 12 items; page 2 returned the remaining 8 without duplicate/missing IDs.
- Province filter `เชียงใหม่`: 5 results. Category filter `อาหารทะเล`: 4 results. Search `ข้าวซอยแม่สาย`: 1 result. Search `สมมติ`: 0 results.
- Health: `200 {"status":"ok"}`. Existing localhost CORS preflight: `204` with the expected allow-origin value.
- All 20 image URLs returned **HTTP 200 / image/jpeg** and were downloaded as actual image files. Representative photographs were visually inspected; no generic placeholder URLs were substituted.
- All imported restaurants initially have `rating: null`, `reviewCount: 0`; external ratings/reviews were not invented or imported.

The supplied FE source already renders `imageUrl` on cards/details and `description` on details. A deployed FE browser session was not tested because no deployed FE URL was supplied. A frontend already configured for this Render API receives the new records after reloading; a CRA mock build must use `REACT_APP_USE_MOCK=false` and the real `REACT_APP_API_URL`, then restart/rebuild. See [the FE handoff](FRONTEND_HANDOFF.md).

The dataset is a snapshot of real businesses and does not guarantee current opening status. Some upstream records lack full street details or coordinates; missing values were not fabricated. External photo-host availability can change, and the FE keeps its normal failure fallback.
