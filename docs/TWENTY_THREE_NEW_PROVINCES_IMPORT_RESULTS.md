# Twenty-three additional provinces: import results

Live public API verification: **2026-11-13 00:38:35 +07 (Asia/Bangkok)**.

The existing ThaiFood API now serves **53 real businesses across 42 provinces**. Twenty-three additions provide one business in each of 23 previously empty provinces, with a different matching licensed photograph and visible attribution. The original thirty restaurant rows are unchanged. [Source facts and image licenses](RESTAURANT_DATA_SOURCES.md) document each record.

| Province | Added business |
| --- | --- |
| จันทบุรี | จันทรโภชนา สาขาเบญจมราชูทิศ |
| ระยอง | Ploy Samet |
| ตราด | Benz Thai Restaurant (Soneva Kiri) |
| ฉะเชิงเทรา | ร้านบะหมี่หยกเป็ดย่าง |
| เพชรบุรี | โอวทึ้ง นายกี๋ (ลอดช่อง) |
| สงขลา | ร้านชมจันทร์ซีฟู้ด เกาะยอ |
| สุพรรณบุรี | ก๋วยเตี๋ยวนายฮวด บางปลาม้า |
| ลพบุรี | บ้านกล้วย&ไข่คาเฟ่ |
| สตูล | The Season by The Cliff Lipe |
| สมุทรปราการ | Bangkok Tree House Restaurant |
| ลำปาง | นิวทะเลเผา |
| ลำพูน | The Terracotta Garden at Lamphun (สวนไม้ไทยพ่อเลี้ยงหมื่น) |
| แม่ฮ่องสอน | ใบเฟิร์น (Fern Restaurant) |
| พะเยา | ครัวออมเงิน |
| นครสวรรค์ | Dairy Queen เซ็นทรัลนครสวรรค์ |
| ขอนแก่น | ภัตตาคารท่าอากาศยานขอนแก่น |
| นครราชสีมา | Rustic Barista |
| บุรีรัมย์ | KLIM Restaurant |
| มุกดาหาร | goodmook |
| เลย | Café Amazon สาขา ปตท.ผานกเค้า |
| ศรีสะเกษ | Café Amazon สำโรง อุทุมพรพิสัย |
| นครนายก | กาแฟกับต้นไม้ |
| นนทบุรี | Tiao Tao Tan |

## Production result

Target: the existing ThaiFood Supabase project `xorsagxmphehnjurquen`, used by `https://thaifood-guide-api-a7c7.onrender.com`. The configured TLS database connection was used without recording credentials in tracked files.

| Check | Before | After addition | After replay |
| --- | ---: | ---: | ---: |
| Restaurants | 30 | 53 | 53 |
| Covered provinces | 19 | 42 | 42 |
| Members | 3 | 3 | 3 |
| Reviews | 0 | 0 | 0 |
| Favorites | 0 | 0 | 0 |
| Reference provinces | 77 | 77 | 77 |
| Categories | 7 | 7 | 7 |

Before writing, the thirty existing restaurant rows were backed up to a private untracked workspace artifact with file permission `0600`. No member, review or favorite contents were exported. A locked transaction rechecked the baseline and all 23 empty target provinces, then ran only the marked `additional_places_23` INSERT block from `seed.sql`. The earlier fictional-row retirement block was not executed.

The original restaurant IDs and all persisted fields retained fingerprint `e0819a8a329494356fe9815ffc6a7e1c`. The replay retained all 53 restaurant IDs/fields with fingerprint `a56ee384a8611c3453c14da1a394913b`. Aggregate member/review/favorite counts stayed unchanged. No application, schema, authentication, CORS, Moviehub or FE repository changes were required; the running API queries PostgreSQL directly, so the additions are already visible without a deployment.

## Verification

- `go vet ./...`: passed.
- `go test -race -count=1 -json ./...`: **73 tests/subtests passed, 0 failed, 0 skipped**, including store and HTTP/API integration against real PostgreSQL 17 in isolated test schemas. Fresh seed/replay checks require 53 restaurants and 42 covered provinces.
- `go build -trimpath`: passed. Formatting and `git diff --check`: passed.
- Before adding records, the updated seed expectation failed against the previous 30-row seed as intended. Final checks used a disposable local PostgreSQL 17 database with at most four test connections per pool; live connection settings were unchanged.
- A read-only review checked the SQL, manifest, source/photo mapping and import safeguards. The original thirty-row seed block is byte-for-byte unchanged apart from its count heading.
- `GET /api/restaurants?limit=50` returns 50 items and page 2 the remaining 3: 53 total across 42 provinces, with 53 distinct nonempty image URLs. All thirty old API rows compare equal field-for-field with the saved public baseline.
- Each new province filter returns exactly its one added business. All 23 detail endpoints match the list and reviewed source manifest, including source/image credits. Floating-point coordinates are compared with a small numeric tolerance.
- Default pagination returns 12, 12, 12, 12 and 5 items with no duplicate or missing IDs. Search for `สมมติ` returns zero.
- Health returns `200 {"status":"ok"}`; reference endpoints return 77 provinces and 7 categories.
- All 23 new image URLs returned **HTTP 200 / image/jpeg**, downloaded as actual images. Each new business has `rating: null` and `reviewCount: 0`; external ratings/reviews were not imported.

The supplied FE source already renders `imageUrl` on cards/details and `description` on details. Reloading an FE configured for this API retrieves the new data. A deployed FE browser session was not tested because no deployed FE URL was supplied. CRA mock builds still need `REACT_APP_USE_MOCK=false` and the real API URL followed by a restart/rebuild; see [the FE handoff](FRONTEND_HANDOFF.md).

These are source snapshots of real businesses; they do not establish current opening status. Missing street details and venue coordinates were not invented. Commons place metadata is explicitly identified where it supports only a locality-level address; camera positions are never stored as venue coordinates.

Changes are prepared on `feat/twenty-three-new-provinces`; this additive data request does not merge that branch or deploy a new application. [Feature branch CI](https://github.com/Sukrit310/thaifood-backend/actions?query=branch%3Afeat%2Ftwenty-three-new-provinces).
