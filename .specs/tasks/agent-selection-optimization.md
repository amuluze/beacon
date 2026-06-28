# Agent Selection Optimization

> Status: active
> Owner: Codex

## Problem

Multi-Agent operations currently depend on each page passing `agent_id` manually. Some monitoring pages do this, but container management and WebSocket log views can fall back to the Server default Agent without a visible target selection.

## Scope

- Provide a shared frontend Agent selection state.
- Send selected Agent consistently through HTTP and WebSocket requests.
- Ensure Agent list API returns stable frontend field names.
- Improve Server-side monitoring query indexes for Agent/time lookups.
- Keep Go test/build usable when the frontend production `dist` directory has not been generated.

## Acceptance

- HTTP API requests include `X-Agent-ID` when a user-selected Agent exists.
- WebSocket log requests include `agent_id` when a user-selected Agent exists.
- Container management pages refresh when the selected Agent changes.
- Agent list JSON exposes `agent_id`, `hostname`, `status`, and `last_seen`.
- Monitoring models have Agent/time composite indexes for trend queries.
- `cd amprobe && go test ./...` is not blocked by a missing frontend `dist` directory.
