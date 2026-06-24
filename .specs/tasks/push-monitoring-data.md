# Task: Push Monitoring Data from Agent to Server

## Status
- [x] Step 1: Common RPC report types (`common/rpc/schema/report.go`)
- [x] Step 2: Server monitoring DB models (`amprobe/service/model/monitor.go`)
- [x] Step 3: Server RPC report service (`amprobe/service/report/`)
- [x] Step 4: Server query path migration (host + container repos)
- [x] Step 5: Server RPC server startup (`amprobe/service/server.go`)
- [x] Step 6: Agent push client + task modification
- [x] Step 7: Agent config & wiring updates
- [x] Step 8: Alarm tasks now read from Server local DB

## Verification
- [x] `go build ./collia/...` passes
- [x] `go build ./amprobe/...` passes
- [x] `go build ./common/...` passes
- [x] `go vet` passes for all modules

## Files Created
- `common/rpc/schema/report.go` — Push report types (MonitorReportArgs/Reply)
- `amprobe/service/model/monitor.go` — Server-side monitoring DB models
- `amprobe/service/report/report.go` — RPC service that stores agent reports
- `amprobe/service/report/rpc_server.go` — RPC server wrapper for report service
- `amprobe/service/report_server.go` — Wire provider for report server
- `collia/service/report/client.go` — Agent-side push client

## Files Modified
- `amprobe/service/model/model.go` — Registered monitoring models
- `amprobe/service/config.go` — Added ReportServer config + ClientNames to TLS
- `amprobe/service/injector.go` — Added ReportServer to injector
- `amprobe/service/wire.go` — Added NewReportServer to wire build
- `amprobe/service/wire_gen.go` — Updated for new signatures
- `amprobe/service/server.go` — Start/stop RPC report server
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
    address: "tcp@amprobe-host:8972"  # Server RPC address
    network: "tcp"                     # or "unix"
    agent_id: "agent-01"              # unique agent identifier
```

### amprobe config.yaml (add to rpc section):
```yaml
rpc:
  report_server:
    network: "tcp"
    address: "0.0.0.0:8972"
    tls:
      enable: false
```
