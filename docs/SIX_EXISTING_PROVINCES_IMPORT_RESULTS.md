# Six additional restaurants in existing provinces

Live verification: **2026-11-13 01:17:57 +07 (Asia/Bangkok)**.

The existing API now serves **59 real businesses across the same 42 provinces**. All six additions are in provinces already represented before this run, with individual matching photographs and visible data/image credits. [Source and license register](RESTAURANT_DATA_SOURCES.md).

| Restaurant | Province |
| --- | --- |
| Anchan Vegetarian Restaurant | เชียงใหม่ |
| WhatCat Cafe | เชียงราย |
| Awesome Coffee Shop | นครราชสีมา |
| ร้านเจ๊ไฝ | กรุงเทพมหานคร |
| Pan and David Restaurant | ชลบุรี |
| โกปี้เตี่ยม (Kopitiam by Wilai) | ภูเก็ต |

| Check | Before | After addition | After replay |
| --- | ---: | ---: | ---: |
| Restaurants | 53 | 59 | 59 |
| Covered provinces | 42 | 42 | 42 |
| Members | 3 | 3 | 3 |
| Reviews | 0 | 0 | 0 |
| Favorites | 1 | 1 | 1 |

A backup of the original restaurant rows was saved locally with mode `0600`; no member/review/favorite contents were exported. Only the marked `additional_existing_places` INSERT block was executed in locked transactions. Earlier fictional-row retirement blocks were not run.
The original 53 restaurant IDs and all database fields retained fingerprint `a56ee384a8611c3453c14da1a394913b`. Replay retained all 59 rows with fingerprint `71d5c2e5276e22d9471da430ea1a5a31`. The province set and member/review/favorite counts stayed unchanged.

- Vet, build, formatting and `git diff --check`: passed.
- `go test -race -count=1 -json ./...`: **73 tests/subtests passed, zero failed/skipped**, including isolated local PostgreSQL 17 integration. The updated 59-row seed expectation first failed against the previous 53-row seed.
- Independent review checked SQL, exact-business source/photo associations and import/preservation safeguards.
- Public list pagination returns 50+9 (maximum-size pages), or 12+12+12+12+11 (defaults), covering all 59 IDs without duplicates. The original 53 API rows compare equal to the saved public baseline.
- All six detail responses match their reviewed records; each selected province filter returns its previous count plus the intended additions. Health is 200/ok and reference endpoints still return 77 provinces/7 categories.
- New businesses have `rating: null` and `reviewCount: 0`; external ratings were not imported.
- All six image URLs returned HTTP 200/image/jpeg and were downloaded as actual images.

The running API reads PostgreSQL directly, so no application deployment was needed. Existing FE code already renders imageUrl/description; an FE configured for this API receives the additions after reloading. No deployed FE browser session was tested. Source snapshots do not guarantee current opening status; absent street/coordinate details remain absent.

Prepared on `feat/six-existing-provinces`; this data request adds no merge or deployment step. [Feature CI](https://github.com/Sukrit310/thaifood-backend/actions?query=branch%3Afeat%2Fsix-existing-provinces).
