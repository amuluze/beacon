> [!IMPORTANT]
> **Repository rename:** this repository has moved from `amuluze/amprobe` to `amuluze/beacon-pro`. GitHub temporarily redirects old links, but existing clones should run `git remote set-url origin git@github.com:amuluze/beacon-pro.git` and update bookmarks, badges, and automation to the new address.

<p align="center">
  <img src="beacon/web/src/assets/images/beacon.png" alt="Beacon Logo" width="128" />
</p>

<h1 align="center">Beacon</h1>

<p align="center">
  <strong>An open-source, lightweight and modern host & Docker monitoring platform</strong><br />
  <sub>开源、轻量、现代化的主机与 Docker 容器监控平台</sub>
</p>

<p align="center">
  <a href="https://github.com/amuluze/beacon/stargazers"><img src="https://img.shields.io/github/stars/amuluze/beacon?style=flat-square" alt="GitHub Stars" /></a>
  <a href="https://github.com/amuluze/beacon/releases"><img src="https://img.shields.io/github/v/release/amuluze/beacon?display_name=tag&sort=semver&style=flat-square" alt="GitHub Release" /></a>
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.25" />
  <img src="https://img.shields.io/badge/Vue-3-42B883?style=flat-square&logo=vuedotjs&logoColor=white" alt="Vue 3" />
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-2EA44F?style=flat-square" alt="MIT License" /></a>
</p>

<p align="center">
  <a href="./README.md">中文</a> · English
</p>

---

<p align="center">
  <a href="#-introduction">Introduction</a> ·
  <a href="#-features">Features</a> ·
  <a href="#-architecture">Architecture</a> ·
  <a href="#-quick-start">Quick Start</a> ·
  <a href="#%EF%B8%8F-tech-stack">Tech Stack</a> ·
  <a href="#-project-structure">Project Structure</a> ·
  <a href="#-documentation">Documentation</a>
</p>

---

## 📖 Introduction

**Beacon** is a Server-Agent platform for host monitoring and Docker container management. It is designed for individual developers and small teams that need one place to observe and operate multiple servers.

- **Beacon Server** provides the Web UI, HTTP API, authentication, authorization, auditing, monitoring storage, and task orchestration.
- **Collia Agent** collects host and Docker metrics, reports monitoring batches over HTTP, and opens a reverse gRPC tunnel for Server-initiated control calls.
- Query and control requests must explicitly select an Agent, preventing implicit fallback or cross-node reads.

Website: [beacon.amuluze.com](https://beacon.amuluze.com) · Repository: [github.com/amuluze/beacon](https://github.com/amuluze/beacon)

## ✨ Features

### 🐳 Docker management

- Inspect Docker versions and runtime status.
- Create, start, stop, restart, delete, and inspect container logs.
- Import, export, delete, and clean dangling images.
- Create, delete, and inspect Docker networks.

### 🖥️ Host monitoring

- Inspect host name, uptime, distribution, kernel, and operating system.
- Track CPU, memory, disk IO, and network IO trends.
- Configure CPU, memory, and disk usage thresholds with email or optionally signed DingTalk group robot Webhook notifications.
- Perform host and container operations through an explicitly selected Agent.

### 🔐 Access control and auditing

- User, role, and API permission management.
- Login, logout, and system operation audit records.
- Production checks for signing keys, Agent join tokens, and install tokens.

### 🔄 Agent lifecycle

- Agent version reporting and online status tracking.
- Collia amd64/arm64 packages are distributed by the Beacon image.
- Remote update, self-update, uninstall, and reverse-tunnel control flows.

## 🏗️ Architecture

```mermaid
flowchart LR
    User["Web UI / API Client"] -->|"HTTP / WebSocket<br/>X-Agent-ID or agent_id"| Beacon["Beacon Server"]
    Beacon --> DB[("Server monitoring DB")]
    Collia["Collia Agent"] -->|"HTTP monitoring report"| Beacon
    Collia -->|"Reverse gRPC tunnel"| Beacon
    Beacon -->|"Control call by Agent ID"| Collia
    Collia --> Runtime["Host OS / Docker Engine"]
```

Beacon keeps three paths distinct:

1. **Monitoring queries** read the selected Agent's data from the Server-side monitoring database.
2. **Monitoring reports** are atomically persisted in batches sent by Collia over HTTP.
3. **Control calls** target an Agent through the reverse gRPC tunnel opened by Collia.

See [System Architecture](./.docs/architecture.md) and [Data Flow](./.docs/concepts/data-flow.md) for the complete boundaries and dependency direction.

## 🚀 Quick Start

### Online installer (recommended)

```bash
bash -c "$(curl -fsSL https://beacon.amuluze.com/release/latest/manager.sh)" -- install BEACON_PUBLIC_BASE_URL=https://beacon.example.com
```

Production mode accepts HTTPS only. First provide a trusted reverse proxy that forwards this HTTPS origin to Beacon's Web port. The installer then configures the ports and generates authentication, data-encryption, and bootstrap-administrator credentials. `BEACON_DATA_ENCRYPTION_KEY` is stored in the mode-`0600` `.env` file and must be retained permanently.

### Run from source

Prerequisites:

- Docker >= 20.10.9 with Docker Compose.
- Go 1.25 for local backend development.
- Node.js, pnpm, and [Task](https://taskfile.dev/) to build Web assets.

```bash
# Clone the renamed repository
git clone https://github.com/amuluze/beacon-pro.git
cd beacon-pro

# Build frontend assets
task beacon-web:install
task beacon-web:build

# Build the Beacon image and start the local Compose stack
docker build -f beacon/Dockerfile -t beacon:latest .
docker compose -f website/docker-compose.yml up -d

# Verify
curl http://127.0.0.1:8000/health
```

The local Compose stack exposes `8000` for HTTP and `17000` for Agent control and is intended for development and verification. Before production deployment, generate unique high-entropy secrets and inject `BEACON_AUTH_SIGNING_KEY`, `BEACON_AGENT_INSTALL_TOKEN`, `BEACON_CONTROL_JOIN_TOKEN`, and `BEACON_DATA_ENCRYPTION_KEY`; a new database also requires `BEACON_BOOTSTRAP_ADMIN_PASSWORD`.

After Beacon starts, install Collia on a target host (replace the address, node number, and token with real values):

```bash
curl -fsSL 'https://<beacon-host>/api/v1/host/install?node=1' | sudo bash -s -- --token=<install-token>
```

### Optional installation reporting

Installation reporting is disabled by default and sends nothing without explicit opt-in. Beacon sends one installation report only after an administrator sets `Enable = true` under `[InstallReport]` in `beacon/configs/config.toml`. The report is sent to the configured `URL`, which defaults to `https://beacon.amuluze.com/api/v1/install/report`.

The payload contains only `install_id`, `image`, `version`, `public_base_url`, `install_dir`, `http_port`, `control_port`, `container_name`, and `hostname`. `install_id` is a random identifier persisted locally. The request does not contain administrator accounts, passwords, authentication tokens, DingTalk Webhooks, SMTP passwords, or monitoring data.

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Web frontend | Vue 3, TypeScript, Vite, Element Plus, Pinia, ECharts |
| Server | Go 1.25, Fiber, GORM, WebSocket |
| Agent | Go, gopsutil, Docker Engine API |
| Control channel | Agent-initiated reverse gRPC tunnel |
| Monitoring channel | HTTP batch reporting and Server-side persistence |
| Storage | SQLite, with other databases available through GORM |
| Deployment | Docker, Docker Compose, Kubernetes |

## 📁 Project Structure

```text
beacon/
├── beacon/                 # Server, Web UI, HTTP/WS API, and tunnel client
│   ├── cmd/beacon/         # Server process entrypoint
│   ├── service/            # Business services, routing, auth, and persistence
│   └── web/                # Vue 3 admin UI
├── collia/                 # Agent collection, host/Docker ops, tunnel service
├── common/                 # Shared schemas, database, and tunnel transport
├── website/                # Product website, installer, and release artifacts
│   ├── compose.yml         # Local dev Compose (NODE_ENV=development)
│   ├── docker-compose.yml  # Production Compose (NODE_ENV=production)
│   └── release/latest/     # Online installer bundle (manager.sh / compose.yaml / SHA256SUMS)
├── .docs/                  # Current implementation documentation
├── .specs/                 # SDD Domain and Task Specs
└── Taskfile.yml            # Workspace command entrypoint
```

## 📚 Documentation

| Document | Contents |
|---|---|
| [Architecture](./.docs/architecture.md) | Module boundaries, runtime paths, and dependency direction |
| [Data Flow](./.docs/concepts/data-flow.md) | Request lifecycle, data ownership, and cross-module flow |
| [Deployment](./.docs/deployment.md) | Environment, build, configuration, and release checks |
| [API Routes](./.docs/api/routes.md) | Current HTTP route and handler index |
| [OpenAPI](./.docs/api/openapi.yml) | Beacon HTTP API contract |
| [Domain Constraints](./.specs/domain/monitoring-platform.md) | Durable monitoring behavior and invariants |

## 🧑‍💻 Development and Contributing

Common commands:

```bash
task beacon:dev
task beacon-web:dev
task collia:amd64

cd beacon && go test ./...
cd collia && go test ./...
cd common && go test ./...
cd beacon/web && pnpm test:run && pnpm ts && pnpm build
```

Issues and pull requests are welcome:

1. Fork [this repository](https://github.com/amuluze/beacon).
2. Create a focused branch and run the relevant checks.
3. Commit a clear, reviewable change.
4. Push the branch and open a Pull Request.

## ☕ Support and Contact

Beacon is maintained in the author's spare time. If it helps you, please give the project a ⭐ or consider buying the author a coffee.

<details>
<summary>Donation and contact QR codes</summary>

<p>
  <img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/beacon/202403171446310.jpg" alt="Donation QR code" width="260" />
  <img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/beacon/202403171449114.jpg" alt="Author WeChat" width="260" />
  <img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/beacon/202403171450306.png" alt="Public account" width="260" />
</p>

</details>

## 📄 License

Beacon is open source under the [MIT License](./LICENSE).

## 🙏 Acknowledgements

Special thanks to [JetBrains](https://www.jetbrains.com/) for supporting open-source development tools.

---

<p align="center">
  <sub>Built with ❤️ and ☕</sub>
</p>
