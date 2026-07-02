# Amprobe

`Amprobe` 项目入口文档，由 `/sdd doc update` 根据当前 workspace 事实重写。

## 30 秒项目摘要

`Amprobe` 是一个 Server-Agent 监控/探测平台：Server 侧提供 Web UI 和 HTTP API，按 Agent 标识选择目标节点；Agent 侧采集主机与 Docker 状态，通过 HTTP 上报监控批次，并通过反向 gRPC tunnel 接收 Server 发起的控制调用。

核心链路：
- Vue/Vite Web 前端发起用户操作或订阅日志、终端等实时通道。
- Fiber Server 接收 HTTP/WebSocket 请求，完成认证授权、参数解析和 Agent 选择。
- 监控查询读取 Server 本地监控表，并按 `X-Agent-ID` 或 `agent_id` 过滤目标 Agent。
- 控制操作通过 `common/rpc/tunnel` 的反向 gRPC tunnel 按 Agent 标识调用目标 Agent。
- Agent 调用 Docker/主机系统 API 完成采集与控制；监控批次通过 HTTP report 入口原子落库，控制结果通过 tunnel reply 或实时流回到 Server。

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
| `amprobe/` | Server control plane: Web/API 接入、认证授权、Agent 生命周期、监控批次落库、目标选择和反向 tunnel client |
| `collia/` | Agent runtime: 主机/容器采集、HTTP 监控上报、Docker 控制和反向 tunnel Service |
| `common/` | shared contract library: 复用 schema、数据库封装、反向 tunnel transport、RPC 参数/返回值和跨模块类型 |
| `deploy/` | supporting project directory |

## Domain Specs

- [.specs/domain/monitoring-platform.md](.specs/domain/monitoring-platform.md)
- [.specs/domain/agent-lifecycle-update.md](.specs/domain/agent-lifecycle-update.md) — Agent 版本上报、远程更新推送、自更新与卸载的领域约束

## 开发命令

```bash
task amprobe:amd64
task amprobe:arm64
task amprobe:bin
task amprobe:build
task amprobe:dev
task amprobe:wire
task web:build
task web:dev
task web:install
task collia:amd64
task collia:arm64
task collia:wire
cd amprobe/web && npm run build
```

## AI Agent 工作流（Spec-Driven Development）

- 新功能或复杂 Bug 先建立 `.specs/tasks/<task-id>.md` 与 `.plans/<task-id>.md`。
- Domain Spec 是最高约束来源，Plan 只描述实施路径。
- 不手工编辑 `.specs/status/*.json`，通过 `/sdd` 命令维护状态。
- `/sdd doc update` 会根据当前项目事实刷新 AGENTS 入口；当 `.docs` 为空时生成默认项目文档，生成或刷新 `.specs/domain/*.md`，并同步 `CLAUDE.md`。

## 技术栈

- reverse gRPC tunnel control channel
- HTTP monitoring report channel
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
- 根目录不是 Go module；Go 验证需进入 `amprobe`、`collia`、`common` 分别执行。
- 监控查询、监控上报和控制调用是三条不同路径；新增接口或修复 Bug 时必须明确所属路径和 Agent 选择语义。
- Agent 选择必须由请求方显式提供 `X-Agent-ID`/`agent_id`；监控查询读路径与控制调用写路径在缺失或格式非法时统一返回错误（`ErrMissingAgentID`/`ErrInvalidAgentID`），禁止回退默认节点或全表查询。`Control.DefaultAgentID` 已废弃。
- 敏感凭据（`Auth.SigningKey`、`Control.JoinToken`、`AgentInstall.Token`）在 `App.Env = production` 下，空值、已知弱默认值或长度不足一律拒绝启动；可用对应环境变量覆盖。
