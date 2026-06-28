# Platform Hardening

> Status: active
> Owner: Codex

## Problem

The monitoring platform has the correct Server-Agent architecture, but several production-facing constraints need stronger implementation coverage:

- Monitoring queries can still fall back to unscoped local table reads when no Agent identity reaches the repository.
- Missing Agent IDs in monitoring reports are recognized by storage but surfaced as generic server errors at the HTTP boundary.
- Agent-side collection can swallow host/network/disk metric errors and report zero-value or partial data as a successful live batch.
- Frontend monitoring views receive freshness metadata but do not mark stale/degraded data.
- Production-facing defaults and duplicate Agent module code need cleanup after the critical correctness fixes.
- Reverse tunnel registration accepts an Agent identity without validating the join token.
- Duplicate Agent IDs can replace an existing tunnel stream.
- Monitoring query responses expose timestamps but do not classify stale data.
- Alarm calculation uses global monitoring tables instead of Agent-scoped windows.
- Some latest-per-group queries depend on database behavior that can return non-latest rows.
- Critical multi-Agent, tunnel, freshness, and error semantics need broader tests.

## Scope

- Harden reverse tunnel registration with optional join-token validation and duplicate/empty Agent ID rejection.
- Wire the Agent join token from generated and runtime configs.
- Require Agent identity at Server-side monitoring query repository boundaries.
- Return input errors for malformed Agent reports instead of generic storage failures.
- Preserve Agent-side collection errors and avoid presenting failed collection as live data.
- Surface stale/degraded freshness state in monitoring UI.
- Remove unused duplicate Agent repository code after confirming it is not wired.
- Add targeted tests for tunnel registration, Agent selection, report storage, freshness, alarms, and latest-per-group queries.
- Add freshness metadata to monitoring query responses without removing existing fields.
- Scope alarm evaluation and cache keys by Agent.
- Replace unsafe latest-per-group SQL with deterministic subqueries.

## Acceptance

- A configured join token is sent by Collia and enforced by the Server tunnel.
- Empty, duplicate, and unauthorized Agent registrations are rejected before lifecycle connect.
- Agent offline, duplicate registration, and unauthorized registration remain distinguishable errors.
- Monitoring responses include stale/degraded metadata derived from collection timestamps.
- Monitoring repositories reject missing Agent identity instead of running unscoped multi-Agent reads.
- Missing `agent_id` reports return a caller-visible input error.
- Failed Agent metric collection is observable and does not silently become zero-value live data.
- Monitoring UI displays stale/degraded state when freshness metadata says data is not live.
- Alarm jobs evaluate each Agent independently and include Agent identity in audit/notification messages.
- Latest disk/container queries return the newest row per resource for the selected Agent.
- `go test ./...` passes for `common`, `amprobe`, and `collia` in an environment that permits local test listeners.
- `pnpm build` passes for `amprobe/web`.
