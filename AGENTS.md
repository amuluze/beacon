# Amprobe

`Amprobe` 项目入口文档，由 `/sdd doc update` 根据当前 workspace 事实重写。

## 30 秒项目摘要

`Amprobe` 是一个 Server-Agent 监控/探测平台：Server 侧提供 Web UI 和 HTTP API，按 Agent 标识选择目标节点，通过 rpcx 调用 Agent；Agent 侧采集主机与 Docker 状态、执行容器控制动作，并把结果落到数据库或返回给前端。

核心链路：
- Vue/Vite Web 前端发起用户操作或订阅日志、终端等实时通道。
- Fiber Server 接收 HTTP/WebSocket 请求，完成认证授权、参数解析和 Agent 选择。
- Server 通过 rpcx client 按 Agent 标识调用目标 Agent 的 RPC Service。
- Agent 读取 GORM 持久化数据或调用 Docker/主机系统 API 完成采集与控制。
- 结果通过共享 schema 或 RPC reply 回到 Server，再转换为前端可展示的响应或实时事件。

## 文档地图

| 文档 | 说明 |
|------|------|
| [.docs/MANIFEST.yml](.docs/MANIFEST.yml) | 文档期望清单，记录 `.docs`、`.specs/domain` 与入口文件的覆盖标准 |
| [.docs/api/routes.md](.docs/api/routes.md) | API route and handler signal index inferred from source registrations |
| [.docs/architecture.md](.docs/architecture.md) | system architecture, runtime flow, module boundaries and dependency direction |
| [.docs/concepts/data-flow.md](.docs/concepts/data-flow.md) | request lifecycle, data ownership boundaries and cross-module flow |
| [.docs/deployment.md](.docs/deployment.md) | environment requirements, build commands, configuration policy and release checks |
| [.docs/modules/amprobe-web.md](.docs/modules/amprobe-web.md) | amprobe-web module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/modules/amprobe.md](.docs/modules/amprobe.md) | amprobe module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/modules/collia.md](.docs/modules/collia.md) | collia module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/modules/common.md](.docs/modules/common.md) | common module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/project-analysis.md](.docs/project-analysis.md) | workspace inventory, documentation health, coverage model and update guidance |
| [docs/server-agent-architecture.md](docs/server-agent-architecture.md) | 项目原有文档；用于补充 `.docs/` 生成视图 |

## 文档健康

- 状态：pass（通过）
- 检查目标：20
- 结构失败：0
- 未纳管信号：0
- 引用错误：0
- 建议：当前无必须处理的文档治理动作；行为变更后继续运行 `/sdd doc update`。

## 关键目录

| Directory | Inferred Role |
|-----------|---------------|
| `.docs/` | implementation documentation |
| `.plans/` | SDD implementation plans |
| `.specs/` | SDD task, status, and domain specs |
| `amprobe/` | Server control plane: Web/API 接入、认证授权、目标选择、RPC client 和运行时协调 |
| `collia/` | Agent runtime: 主机/容器采集、Docker 控制、GORM 本地状态和 rpcx Service |
| `common/` | shared contract library: 复用 schema、数据库封装、RPC 参数/返回值和跨模块类型 |
| `deploy/` | supporting project directory |
| `docs/` | project documentation |
| `installer/` | supporting project directory |

## Domain Specs

- [.specs/domain/monitoring-platform.md](.specs/domain/monitoring-platform.md)

## 开发命令

```bash
cd amprobe && make amd64
cd amprobe && make arm64
cd amprobe && make bin
cd amprobe && make build
cd amprobe && make dev
cd amprobe && make wire
cd amprobe/web && make build
cd amprobe/web && make dev
cd amprobe/web && make install
cd collia && make amd64
cd collia && make arm64
cd collia && make installer
cd collia && make wire
cd amprobe/web && npm run build
```

## AI Agent 工作流（Spec-Driven Development）

- 新功能或复杂 Bug 先建立 `.specs/tasks/<task-id>.md` 与 `.plans/<task-id>.md`。
- Domain Spec 是最高约束来源，Plan 只描述实施路径。
- 不手工编辑 `.specs/status/*.json`，通过 `/sdd` 命令维护状态。
- `/sdd doc update` 会根据当前项目事实刷新 AGENTS 入口；当 `.docs` 为空时生成默认项目文档，生成或刷新 `.specs/domain/*.md`，并同步 `CLAUDE.md`。

## 技术栈

- rpcx Server-Agent RPC
- Docker/host operation boundary
- Vue/Vite frontend
- Fiber/HTTP API
- Server-Agent monitoring/probing domain
- GORM persistence
- WebSocket realtime channel
- Go modules detected from `go.mod`/`go.work`
- Node.js/Vue/TypeScript package metadata detected from `package.json`
- Markdown documentation under `.docs/`
- SDD domain specs under `.specs/domain/`

## 关键约束

- `.specs/domain/` 描述长期约束空间，只写必须满足的行为、不变量、状态和错误语义，不写实现方案。
- `.docs/` 描述当前实现事实，结构和内容必须能从源码、配置、构建脚本或现有文档追溯。
- `AGENTS.md` 是 AI 协作入口的 SSOT；`CLAUDE.md` 是同步副本，内容必须一致。
- 新增模块、命令、路由、配置或领域概念后，重新运行 `/sdd doc update` 并检查 `.docs/MANIFEST.yml`。
- 不在文档中写入真实密钥、Token、内部凭据或不可公开的环境值。
