> [!IMPORTANT]
> **仓库更名说明：** 本仓库已由 `amuluze/amprobe` 更名为 `amuluze/beacon`。GitHub 会暂时重定向旧链接，但请在已有本地仓库中执行 `git remote set-url origin git@github.com:amuluze/beacon.git`，并将收藏、徽章与自动化配置更新为新地址。

<p align="center">
  <img src="beacon/web/src/assets/images/beacon.png" alt="Beacon Logo" width="128" />
</p>

<h1 align="center">Beacon</h1>

<p align="center">
  <strong>开源、轻量、现代化的主机与 Docker 容器监控平台</strong><br />
  <sub>An open-source, lightweight and modern host & Docker monitoring platform</sub>
</p>

<p align="center">
  <a href="https://github.com/amuluze/beacon/stargazers"><img src="https://img.shields.io/github/stars/amuluze/beacon?style=flat-square" alt="GitHub Stars" /></a>
  <a href="https://github.com/amuluze/beacon/releases"><img src="https://img.shields.io/github/v/release/amuluze/beacon?display_name=tag&sort=semver&style=flat-square" alt="GitHub Release" /></a>
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.25" />
  <img src="https://img.shields.io/badge/Vue-3-42B883?style=flat-square&logo=vuedotjs&logoColor=white" alt="Vue 3" />
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-2EA44F?style=flat-square" alt="MIT License" /></a>
</p>

<p align="center">
  中文 · <a href="./README.en.md">English</a>
</p>

---

<p align="center">
  <a href="#-项目简介">项目简介</a> ·
  <a href="#-功能特性">功能特性</a> ·
  <a href="#-系统架构">系统架构</a> ·
  <a href="#-快速开始">快速开始</a> ·
  <a href="#%EF%B8%8F-技术栈">技术栈</a> ·
  <a href="#-项目结构">项目结构</a> ·
  <a href="#-项目文档">项目文档</a>
</p>

---

## 📖 项目简介

**Beacon** 是一个采用 Server-Agent 架构的主机监控与 Docker 容器管理平台，面向需要统一观察、管理多台服务器的个人开发者与小型团队。

- **Beacon Server** 提供 Web UI、HTTP API、认证授权、审计、监控数据存储与任务编排。
- **Collia Agent** 采集主机和 Docker 指标，通过 HTTP 上报监控批次，并主动建立反向 gRPC tunnel 接收 Server 控制调用。
- 查询与控制请求必须显式指定 Agent，避免默认节点回退或跨节点读取。

官网：[help.beacon.amuluze.com](https://help.beacon.amuluze.com) · 仓库：[github.com/amuluze/beacon](https://github.com/amuluze/beacon)

## 🖼️ 产品截图

![Beacon overview](website/web/public/images/overview.png)

<details>
<summary>查看更多界面</summary>

| 容器监控 | 主机监控 |
|---|---|
| ![Container monitoring](website/web/public/images/container_monitor.png) | ![Host monitoring](website/web/public/images/host_monitor.png) |

</details>

## ✨ 功能特性

### 🐳 Docker 管理

- 查看 Docker 版本与运行状态。
- 管理容器的创建、启动、停止、重启、删除与日志。
- 管理镜像导入、导出、删除和虚悬镜像清理。
- 创建、删除并查看 Docker 网络状态。

### 🖥️ 主机监控

- 查看主机名、启动时间、发行版、内核与系统类型。
- 观察 CPU、内存、磁盘 IO 与网络 IO 趋势。
- 通过明确选择的 Agent 执行主机与容器控制操作。

### 🔐 权限与审计

- 用户、角色与接口权限管理。
- 登录、登出和系统操作审计。
- 生产模式下校验签名密钥、Agent 加入令牌与安装令牌。

### 🔄 Agent 生命周期

- Agent 版本上报与在线状态跟踪。
- Collia amd64/arm64 安装包由 Beacon 镜像统一提供。
- 支持远程更新、自更新、卸载和反向 tunnel 控制。

## 🏗️ 系统架构

```mermaid
flowchart LR
    User["Web UI / API Client"] -->|"HTTP / WebSocket<br/>X-Agent-ID or agent_id"| Beacon["Beacon Server"]
    Beacon --> DB[("Server monitoring DB")]
    Collia["Collia Agent"] -->|"HTTP monitoring report"| Beacon
    Collia -->|"Reverse gRPC tunnel"| Beacon
    Beacon -->|"Control call by Agent ID"| Collia
    Collia --> Runtime["Host OS / Docker Engine"]
```

Beacon 将三条路径明确分离：

1. **监控查询**：Web 请求从 Server 本地监控表读取指定 Agent 的数据。
2. **监控上报**：Collia 通过 HTTP report 接口按批次原子写入 Server。
3. **控制调用**：Beacon 通过 Collia 主动建立的反向 gRPC tunnel 调用目标 Agent。

更完整的边界、依赖方向和数据流见 [架构文档](./.docs/architecture.md) 与 [数据流文档](./.docs/concepts/data-flow.md)。

## 🚀 快速开始

### 在线安装（推荐）

```bash
curl -fsSL https://help.beacon.amuluze.com/download/install.sh | sh
```

安装脚本会引导配置 Web 端口、Agent 控制端口与安全凭据。

### 从源码启动

环境要求：

- Docker >= 20.10.9，并安装 Docker Compose。
- Go 1.25（本地开发后端时需要）。
- Node.js、pnpm 与 [Task](https://taskfile.dev/)（构建 Web 资源时需要）。

```bash
# 克隆新仓库
git clone https://github.com/amuluze/beacon.git
cd beacon

# 构建前端静态资源
task beacon-web:install
task beacon-web:build

# 构建 Beacon 镜像并启动本地 Compose
docker build -f beacon/Dockerfile -t beacon:latest .
docker compose -f deploy/docker-compose.yml up -d

# 验证服务
curl http://127.0.0.1:8000/health
```

本地 Compose 暴露 `8000`（HTTP）与 `17000`（Agent control）端口，仅用于开发/验证。生产部署前请生成独立的高强度密钥，并使用 `BEACON_AUTH_SIGNING_KEY`、`BEACON_AGENT_INSTALL_TOKEN`、`BEACON_CONTROL_JOIN_TOKEN` 注入。

Beacon 启动后，可从目标主机安装 Collia Agent（将地址、节点编号和 Token 替换为实际值）：

```bash
curl -kfsSL 'http://<beacon-host>:8000/api/v1/host/install?node=1' | sudo bash -s -- --token=<install-token>
```

## 🛠️ 技术栈

| 层级 | 技术 |
|---|---|
| Web 前端 | Vue 3、TypeScript、Vite、Element Plus、Pinia、ECharts |
| Server | Go 1.25、Fiber、GORM、WebSocket |
| Agent | Go、gopsutil、Docker Engine API |
| 控制通道 | Agent 主动连接的反向 gRPC tunnel |
| 监控通道 | HTTP 批次上报与 Server 本地持久化 |
| 数据存储 | SQLite，支持通过 GORM 扩展其他数据库 |
| 部署 | Docker、Docker Compose、Kubernetes |

## 📁 项目结构

```text
beacon/
├── beacon/                 # Server、Web UI、HTTP/WS API 与 tunnel client
│   ├── cmd/beacon/         # Server 进程入口
│   ├── service/            # 业务服务、路由、认证与数据访问
│   └── web/                # Vue 3 管理端
├── collia/                 # Agent 采集、Docker/主机操作与 tunnel service
├── common/                 # 共享 schema、数据库与 reverse tunnel transport
├── website/                # Beacon 官网与安装脚本服务
├── deploy/                 # Docker Compose 与 Kubernetes 部署文件
├── .docs/                  # 当前实现文档
├── .specs/                 # SDD Domain / Task Specs
└── Taskfile.yml            # 工作区开发命令入口
```

## 📚 项目文档

| 文档 | 内容 |
|---|---|
| [架构设计](./.docs/architecture.md) | 模块边界、核心运行链路与依赖方向 |
| [数据流](./.docs/concepts/data-flow.md) | 请求生命周期、数据归属与跨模块流转 |
| [部署指南](./.docs/deployment.md) | 环境、构建、配置与发布检查 |
| [API 路由](./.docs/api/routes.md) | 当前 HTTP 路由与处理器索引 |
| [OpenAPI](./.docs/api/openapi.yml) | Beacon HTTP API 契约 |
| [领域约束](./.specs/domain/monitoring-platform.md) | 监控平台长期行为与不变量 |

## 🧑‍💻 开发与贡献

常用命令：

```bash
task beacon:dev
task beacon-web:dev
task collia:amd64

cd beacon && go test ./...
cd collia && go test ./...
cd common && go test ./...
cd beacon/web && pnpm test:run && pnpm ts && pnpm build
```

欢迎提交 Issue 和 Pull Request：

1. Fork [本仓库](https://github.com/amuluze/beacon)。
2. 创建功能分支并完成最小验证。
3. 提交清晰、可审查的变更。
4. 推送分支并创建 Pull Request。

## ☕ 支持项目与联系作者

Beacon 由作者利用业余时间持续维护。如果项目对你有帮助，欢迎给仓库点一个 ⭐，也可以请作者喝杯咖啡。

<details>
<summary>展开赞赏码与联系方式</summary>

<p>
  <img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/beacon/202403171446310.jpg" alt="赞赏码" width="260" />
  <img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/beacon/202403171449114.jpg" alt="作者微信" width="260" />
  <img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/beacon/202403171450306.png" alt="公众号" width="260" />
</p>

</details>

## 📄 License

Beacon 基于 [MIT License](./LICENSE) 开源。

## 🙏 鸣谢

特别感谢 [JetBrains](https://www.jetbrains.com/) 为开源项目提供开发工具支持。

---

<p align="center">
  <sub>用 ❤️ 与 ☕ 打造 · Built with ❤️ and ☕</sub>
</p>
