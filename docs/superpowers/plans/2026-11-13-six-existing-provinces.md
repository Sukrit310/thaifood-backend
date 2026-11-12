# Six additional businesses within existing coverage

User request: add six real restaurants within the existing 42 covered provinces. Preserve all 53 existing restaurants and all user data; final seed/live counts are 59 restaurants across the same 42 provinces.

1. Verify six distinct businesses and matched freely reusable photographs against TAT/government/operator/Commons evidence. Record sources, author and license; do not invent coordinates, address components, ratings or opening status.
2. Update existing seed replay/API tests to require 59/42 and pagination 50+9. Verify the old 53-row seed fails the updated count check; append only the six reviewed rows.
3. Run vet, full race tests with isolated local PostgreSQL 17, build, formatting and independent review. Preserve the user's `.env.example` modification.
4. Back up current restaurant rows with mode 0600. Apply only the new INSERT block in locked transactions. Require unchanged old fields, identical province set, exactly one addition per selected province, unchanged member/review/favorite counts and idempotent replay.
5. Check the live list, six details and province filters, existing 53 API rows, pagination 50+9/default12+12+12+12+11, health/reference data and all six image URLs. Record actual results, commit selected task files, push the feature branch and confirm CI.

No API/schema/FE/CORS/deployment change or new merge step is part of this data request. Credentials, backups and raw evidence stay outside tracked files.
