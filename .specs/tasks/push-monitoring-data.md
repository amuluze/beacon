# Task: Push Monitoring Data from Agent to Server

## Status
- [x] Step 1: Common report schema (`common/rpc/schema/report.go`)
- [x] Step 2: Server monitoring DB models (`amprobe/service/model/monitor.go`)
- [x] Step 3: Server HTTP report service (`amprobe/service/report/`)
- [x] Step 4: Server query path migration (host + container repos)
- [x] Step 5: Server HTTP route registration (`POST /api/v1/host/report`)
- [x] Step 6: Agent push client + task modification
- [x] Step 7: Agent config & wiring updates
- [x] Step 8: Alarm tasks now read from Server local DB
- [x] Step 9: Report storage rejects missing `agent_id` and stores a batch as one consistency unit

## Verification
- [x] `cd collia && go test ./...` passes
- [x] `cd amprobe && go test ./...` passes
- [x] `cd common && go test ./...` passes
- [x] `cd collia && go build ./...` passes
- [x] `cd amprobe && go build ./...` passes
- [x] `cd common && go build ./...` passes

## Files Created
- `common/rpc/schema/report.go` — Push report types (MonitorReportArgs/Reply)
- `amprobe/service/model/monitor.go` — Server-side monitoring DB models
- `amprobe/service/report/report.go` — HTTP report service that stores agent reports
- `collia/service/report/client.go` — Agent-side push client

## Files Modified
- `amprobe/service/model/model.go` — Registered monitoring models
- `amprobe/service/router.go` — Registered `POST /api/v1/host/report`
- `amprobe/service/host/repository/host.go` — Monitoring queries from local DB
- `amprobe/service/container/repository/container.go` — Monitoring queries from local DB
- `amprobe/service/task/task.go` — Alarm checks from local DB instead of RPC
- `amprobe/service/task.go` — Updated NewTimedTask signature
- `collia/service/task/task.go` — Report-based collection (no DB writes)
- `collia/service/task/host.go` — Collects data, returns report structs
- `collia/service/task/container.go` — Collects data, returns report structs
- `collia/service/task.go` — Push to Server instead of local DB
- `collia/service/config.go` — Added report config section
- `collia/service/wire_gen.go` — Updated for new TimedTask signature

## Configuration Changes Required

### collia config.yaml (add to task section):
```yaml
task:
  report:
    url: "http://amprobe-host:8000/api/v1/host/report"
    token: ""
    agent_id: "agent-01"              # unique agent identifier
```

### amprobe config.toml

`[AgentInstall] Token` controls report/install token semantics. The report route is served by the existing Fiber HTTP server.
