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

### After (Push model)
```
collia (collect) → RPC push → amprobe RPC Server → central DB
Frontend → HTTP API → amprobe → central DB
```

## Implementation Steps

### Step 1: Common RPC Types
- File: `common/rpc/schema/report.go` — push report request/reply types

### Step 2: Server — Monitoring DB Models
- File: `amprobe/service/model/monitor.go` — GORM models matching collia's current models
- Update: `amprobe/service/model/model.go` — register new models

### Step 3: Server — RPC Report Service
- File: `amprobe/service/report/report.go` — receives push data, stores in DB
- File: `amprobe/service/report/rpc_server.go` — rpcx server registration and lifecycle

### Step 4: Server — Query Path Migration
- Modify: `amprobe/service/host/repository/host.go` — read from local DB, not RPC
- Modify: `amprobe/service/container/repository/container.go` — read from local DB, not RPC
- Update: `amprobe/service/task/task.go` — alarm checks read from local DB

### Step 5: Server — RPC Server Startup
- Modify: `amprobe/service/server.go` — start RPC server alongside HTTP server
- Update: `amprobe/service/config.go` — add RPC server config (listen address)

### Step 6: Agent — Push Client
- File: `collia/service/report/client.go` — RPC client to push data to Server
- Modify: `collia/service/task/host.go` — push instead of local store
- Modify: `collia/service/task/container.go` — push instead of local store
- Modify: `collia/service/task.go` — remove local DB write for monitoring, add push

### Step 7: Agent — Config & Wiring
- Modify: `collia/service/config.go` — add Server RPC address config
- Update: `collia/service/wire.go` / `wire_gen.go` — wire report client

## Key Design Decisions

1. **RPC Direction**: Agent pushes to Server (new RPC server on Server side)
2. **DB Schema**: Same table structure as current Agent models, just on Server
3. **Data Flow**: Agent collects → pushes immediately → Server stores → Frontend queries
4. **Agent DB**: Remove `s_host`, `s_cpu`, `s_memory`, `s_disk`, `s_net`, `s_container`, `s_docker`, `s_image`, `s_network` tables from Agent
5. **System Operations**: Keep file/DNS/time/reboot/shutdown operations on Agent (they require local system access)
