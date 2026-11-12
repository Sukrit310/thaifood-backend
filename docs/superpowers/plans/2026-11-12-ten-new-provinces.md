# Add ten real restaurants in previously empty provinces

User request: add ten real businesses in provinces without existing restaurant data. Use one per province and matched, licensed photographs. This is an additive data change; no API/schema/FE redesign is needed.

1. Record live coverage (20 restaurants / 9 provinces) and select ten distinct empty provinces.
2. Verify four factual TAT records plus six from government/Commons place sources; match each photo to the exact business and record author/license/source in its description. Do not invent missing street details, coordinates, ratings or current opening status.
3. Append an INSERT-only, source-deduplicated block to seed.sql. Update existing fresh/replay seed and API checks to require 30 restaurants and 19 provinces. Keep the previous 20 records and retention logic intact.
4. Run vet, race tests including real PostgreSQL integration in temporary schemas, build, formatting and diff checks. Review the final mapping; exclude the separate JR Tham Krasae entry from the photographed Tham Krasae business.
5. Back up restaurant rows with restrictive file permissions, then import only the additive block in a transaction. Abort on changed baseline, occupied target province, altered old records or changed member/review/favorite counts. Replay and compare all restaurant fields.
6. Verify the public list, details, ten province filters, pagination, reference data and all ten image URLs. Document actual evidence, commit the selected task files and push the feature branch. This data request does not add a new merge or deployment step.

Provider totals afterward: 23 TAT, 1 operator, 3 other government, 3 Commons place metadata. Missing venue coordinates remain null. Keep credentials and private backup/evidence files outside tracked changes; preserve the user's existing .env.example modification.
