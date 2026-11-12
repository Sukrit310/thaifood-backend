# Ten additional provinces: import results

Latest additive data run: [23 more provinces and current 53-restaurant verification](TWENTY_THREE_NEW_PROVINCES_IMPORT_RESULTS.md). Earlier results below remain historical.

Live public API verification: **2026-11-13 00:08:29 +07 (Asia/Bangkok)**.

The existing ThaiFood API now serves **30 real businesses across 19 provinces**. Ten additions provide one business in each of ten previously empty provinces, with a different matching licensed photograph and visible attribution. The original twenty restaurant rows are unchanged. [Source facts and image licenses](RESTAURANT_DATA_SOURCES.md) document each record.

| Province | Added business |
| --- | --- |
| สิงห์บุรี | รำพึง |
| กาญจนบุรี | ร้านอาหารถ้ำกระแซ |
| ชุมพร | แม่เฉลียว |
| พังงา | ร้านอาหารปันหยี |
| สุราษฎร์ธานี | The Mountain Grand View |
| ตรัง | ขนมจีบสังขยาป้าพิณ |
| พระนครศรีอยุธยา | ก๋วยเตี๋ยวปากหม้อ พนมสารคาม บางปะอิน |
| ปทุมธานี | ร้านกาแฟชายทุ่ง |
| พิษณุโลก | Café Amazon พลายชุมพล |
| สระแก้ว | Café Amazon สตาร์พลาซ่า อรัญประเทศ |

## Production result

Target: the existing ThaiFood Supabase project `xorsagxmphehnjurquen`, used by `https://thaifood-guide-api-a7c7.onrender.com`. The configured TLS database connection was used without recording credentials in tracked files.

| Check | Before | After addition | After replay |
| --- | ---: | ---: | ---: |
| Restaurants | 20 | 30 | 30 |
| Covered provinces | 9 | 19 | 19 |
| Members | 2 | 2 | 2 |
| Reviews | 0 | 0 | 0 |
| Favorites | 0 | 0 | 0 |
| Reference provinces | 77 | 77 | 77 |
| Categories | 7 | 7 | 7 |

Before writing, the twenty existing restaurant rows were backed up to a private untracked workspace artifact with file permission `0600`. No member, review or favorite contents were exported. A locked transaction rechecked the baseline and all ten empty target provinces, then ran only the marked `additional_places` INSERT block from `seed.sql`. The earlier fictional-row retirement block was not executed.

The original restaurant IDs and all persisted fields retained fingerprint `4f7a04f83dadb6de0df5567e5d39f43f`. The replay retained all thirty restaurant IDs/fields with fingerprint `e0819a8a329494356fe9815ffc6a7e1c`. Aggregate member/review/favorite counts stayed unchanged. No application, schema, authentication, CORS, Moviehub or FE repository changes were required; the running API queries PostgreSQL directly, so the additions are already visible without a deployment.

## Verification

- `go vet ./...`: passed.
- `go test -race -count=1 -json ./...`: **73 tests/subtests passed, 0 failed, 0 skipped**, including store and HTTP/API integration against real PostgreSQL 17 in isolated test schemas. Fresh seed/replay checks require 30 restaurants and 19 covered provinces.
- `go build -trimpath`: passed. Formatting and `git diff --check`: passed.
- The first test run against the hosted connection passed 72 tests, then failed while acquiring a connection (`SQLSTATE 53300`, reserved connection slots). The complete successful run used a disposable local PostgreSQL 17 database with at most four test connections per pool; live connection settings were unchanged.
- A read-only review checked the restored SQL, manifest, source/photo mapping and import safeguards. The original twenty-row seed block is byte-for-byte unchanged apart from its count heading.
- `GET /api/restaurants?limit=50`: 30 items across 19 provinces, with 30 distinct nonempty image URLs. All twenty old API rows compare equal field-for-field with the saved public baseline.
- Each new province filter returns exactly its one added business. All ten detail endpoints match the list and reviewed source manifest, including source/image credits. Floating-point coordinates are compared with a small numeric tolerance.
- Default pagination returns 12, 12 and 6 items with no duplicate or missing IDs. Search for `สมมติ` returns zero.
- Health returns `200 {"status":"ok"}`; reference endpoints return 77 provinces and 7 categories.
- All ten new image URLs returned **HTTP 200 / image/jpeg**, downloaded as actual images. Each new business has `rating: null` and `reviewCount: 0`; external ratings/reviews were not imported.

The supplied FE source already renders `imageUrl` on cards/details and `description` on details. Reloading an FE configured for this API retrieves the new data. A deployed FE browser session was not tested because no deployed FE URL was supplied. CRA mock builds still need `REACT_APP_USE_MOCK=false` and the real API URL followed by a restart/rebuild; see [the FE handoff](FRONTEND_HANDOFF.md).

These are source snapshots of real businesses; they do not establish current opening status. Missing street details and venue coordinates were not invented. ร้านอาหารถ้ำกระแซ is the photographed business, separate from TAT’s เจ.อาร์. ถ้ำกระแซ entry. Commons locality evidence is explicitly identified for that restaurant and the two Café Amazon branches.

Changes are prepared on `feat/ten-new-provinces`; this additive data request does not merge that branch or deploy a new application. [Feature branch CI](https://github.com/Sukrit310/thaifood-backend/actions?query=branch%3Afeat%2Ften-new-provinces).
