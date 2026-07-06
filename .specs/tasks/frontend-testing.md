# Frontend Testing

> Status: active
> Owner: Codex

## Problem

The amprobe-web frontend has zero automated tests. The `agent` store, Axios interceptors, and WebSocket handlers are critical to multi-Agent correctness but have no regression coverage. Manual QA is the only verification path, which slows iteration and increases the risk of breaking Agent selection, token refresh, or request routing.

## Scope

- Add Vitest + jsdom + @vue/test-utils to the frontend toolchain.
- Cover `useAgentStore` (agent selection, fallback, persistence).
- Cover Axios request interceptor (`X-Agent-ID` injection, token attachment).
- Cover Axios response interceptor (401 refresh flow, 403/500 error handling).
- Add `pnpm test` to `Taskfile.yml` and CI steps.

## Acceptance

- `pnpm test` runs in `amprobe/web` and exits 0.
- `useAgentStore` tests pass: setAgents, setCurrentAgent, empty-list fallback.
- Axios interceptor tests pass: `X-Agent-ID` present when agent selected, token present when logged in.
- Response interceptor tests pass: 401 triggers token refresh queue, 403 shows permission warning.
- Tests do not require a running backend server.
