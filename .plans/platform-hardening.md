# Platform Hardening Plan

## Steps

1. P0: Require Agent identity for local monitoring queries and return clear input errors for malformed reports.
2. P0: Preserve Collia collection errors, harden psutil network collection against panics, and avoid silent zero-value live batches.
3. P0: Display stale/degraded freshness state in monitoring UI.
4. P1: Keep reverse tunnel join-token validation, empty Agent ID rejection, duplicate Agent ID rejection, and config wiring covered by tests.
5. P1: Remove unused duplicate Agent repository code and tighten production-facing defaults where this can be done without breaking dev/test flows.
6. Run Go module tests and frontend build; record any environment-specific caveats.
