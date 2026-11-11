# Postman verification

Import `ThaiFood_Guide.postman_collection.json` and `local.postman_environment.json`. Set `baseUrl`, `adminEmail` and `adminPassword` in a private environment. Run requests in order against a test environment; the collection creates its own members and restaurant, then deletes that restaurant. Do not commit real credentials or exported run reports.

The default `allowEdgeSqlProbeBlock=false` keeps all Contract checks strict. The final SQL-injection probe must return HTTP 200, JSON, and an empty literal-search result when it reaches Go. Direct Go → Supabase integration tests and the Supabase-backed Postman run verify this behavior.

For verification through Render, explicitly set `allowEdgeSqlProbeBlock=true` if its upstream WAF blocks this probe. This exception applies only to request **64 SQL metacharacters remain data**, an HTTPS `*.onrender.com` base URL, HTTP 403, `Server: cloudflare`, a nonempty `CF-Ray`, HTML content type, and the exact WAF block message. Every condition is required. Generic 403/HTML responses still fail, and every other request retains the strict JSON/error-envelope checks.

The first live Render run passed the other 63 requests but had four assertion failures on this final request because the upstream WAF returned its HTML block page before the Go API. The explicit option records that security-layer result separately; it does not prove that this request reached Go or replace the direct database-backed SQL-parameterization tests. It does not change the API Contract or any backend behavior.
