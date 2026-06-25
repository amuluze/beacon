# Plan: Push Monitoring Data from Agent to Server

## Problem

Currently, collia (Agent) collects host/container monitoring data and stores it in local SQLite. amprobe (Server) queries the data via RPC when the frontend requests it. This means:

- Data is scattered across Agent nodes
- Each Agent needs its own database
- Historical data retention is limited by Agent disk space
- Alarm checks require RPC calls to each Agent

## Goal

Move monitoring data storage from Agent to Server:

1. **Agent (collia)**: Collects data, pushes to Server in real-time after collection
2. **Server (amprobe)**: Receives data from all Agents, stores centrally, serves to frontend

## Architecture Change

### Before (Pull model)
```
Frontend → HTTP API → amprobe RPC Client → collia RPC Server → SQLite
```

### After (Push model, current implementation)
```
collia (collect) → HTTP POST /api/v1/host/report → amprobe → central DB
Frontend → HTTP API → amprobe → central DB
```

## Implementation Steps

### Step 1: Common RPC Types
- File: `common/rpc/schema/report.go` — push report request/reply types

### Step 2: Server — Monitoring DB Models
- File: `amprobe/service/model/monitor.go` — GORM models matching collia's current models
- Update: `amprobe/service/model/model.go` — register new models

### Step 3: Server — HTTP Report Service
- File: `amprobe/service/report/report.go` — receives push data, validates Agent identity, stores batch in DB
- Route: `POST /api/v1/host/report` — registered in Fiber router

### Step 4: Server — Query Path Migration
- Modify: `amprobe/service/host/repository/host.go` — read from local DB, not RPC
- Modify: `amprobe/service/container/repository/container.go` — read from local DB, not RPC
- Update: `amprobe/service/task/task.go` — alarm checks read from local DB

### Step 5: Server — HTTP Route Startup
- Modify: `amprobe/service/router.go` — register report route on existing HTTP server
- Config: reuse existing HTTP server and Agent install/report token semantics

### Step 6: Agent — Push Client
- File: `collia/service/report/client.go` — RPC client to push data to Server
- Modify: `collia/service/task/host.go` — push instead of local store
- Modify: `collia/service/task/container.go` — push instead of local store
- Modify: `collia/service/task.go` — remove local DB write for monitoring, add push

### Step 7: Agent — Config & Wiring
- Modify: `collia/service/config.go` — add Server RPC address config
- Update: `collia/service/wire.go` / `wire_gen.go` — wire report client

## Key Design Decisions

1. **Report Direction**: Agent pushes monitoring data to Server via HTTP POST
2. **DB Schema**: Same table structure as current Agent models, just on Server
3. **Data Flow**: Agent collects → pushes immediately → Server stores → Frontend queries
4. **Agent DB**: Remove `s_host`, `s_cpu`, `s_memory`, `s_disk`, `s_net`, `s_container`, `s_docker`, `s_image`, `s_network` tables from Agent
5. **System Operations**: Keep file/DNS/time/reboot/shutdown operations on Agent (they require local system access)
6. **Batch Consistency**: Server accepts or rejects an Agent report batch as one consistency unit; missing `agent_id` is rejected.
