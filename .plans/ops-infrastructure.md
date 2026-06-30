# Ops Infrastructure Plan

## Steps

1. **Health probes** (P1)
   - Add `amprobe/service/health/api/health.go` with `Probe` struct and `Liveness`/`Readiness` handlers.
   - Register `GET /health` and `GET /ready` in `amprobe/service/router.go` before auth middleware.
   - Wire `DBHealthy` and `TunnelHealthy` checks when interfaces are available.

2. **Docker Compose** (P2)
   - Create `docker-compose.yml` in project root with `amprobe` service, named volumes, config mount, and healthcheck.
   - Create `deploy/docker-compose.yml` as a reference for production-like setups.

3. **Rate limiting** (P1)
   - Import `gofiber/fiber/v2/middleware/limiter` in `router.go`.
   - Apply limiter to `/api/v1/auth/login` (max 10 req/min per IP) and `/api/v1/host/report` (max 60 req/min per Agent).

4. **Dockerfile alignment** (P0)
   - Update `amprobe/Dockerfile` builder images from `golang:1.21` to `golang:1.25` to match workspace.

5. **Validation**
   - `cd amprobe && go test ./...` passes.
   - `docker-compose up` starts and `/health` responds.
