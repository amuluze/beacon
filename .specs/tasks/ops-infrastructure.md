# Ops Infrastructure

> Status: active
> Owner: Codex

## Problem

The project lacks production-facing operational infrastructure: no health/readiness probes for container orchestration, no rate limiting for API abuse prevention, no docker-compose for local development, and no CI/CD pipeline configuration. These gaps make production deployment risky and increase the operational burden.

## Scope

- Add `/health` (liveness) and `/ready` (readiness) HTTP endpoints to the Server.
- Add rate limiting middleware for login and Agent report endpoints.
- Provide `docker-compose.yml` for one-command local Server startup.
- Wire health probes into dependency checks (DB connectivity, tunnel state) once interfaces are available.
- Update Dockerfile base image to match Go workspace version.

## Acceptance

- `GET /health` returns HTTP 200 with `{"status":"alive","uptime":...}`.
- `GET /ready` returns HTTP 200 when all dependencies are healthy; HTTP 503 when DB or tunnel is down.
- `docker-compose up` starts the Server with persistent SQLite volume and mapped configs.
- Rate limiting rejects excessive requests to `/api/v1/auth/login` and `/api/v1/host/report` with HTTP 429.
- Health endpoints are not protected by auth/Casbin middleware.
- `go test ./...` passes for `amprobe`, `collia`, and `common`.
