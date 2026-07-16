# Beacon

`Beacon` 项目入口文档，由 `/doc update` 根据当前 workspace 事实重写。

## 30 秒项目摘要

`Beacon` 是包含 Go 后端和前端包的应用，后端提供运行时能力，前端通过 API 与后端协作。

核心链路：
- 前端页面通过 API client 发起请求。
- 后端 HTTP 入口解析请求并调用应用服务。
- 服务层完成状态读取、持久化或外部副作用后返回响应。

## 文档地图

| 文档 | 说明 |
|------|------|
| [.docs/MANIFEST.yml](.docs/MANIFEST.yml) | 文档期望清单，记录 `.docs`、`.specs/domain` 与入口文件的覆盖标准 |
| [.docs/api/routes.md](.docs/api/routes.md) | API route and handler signal index inferred from source registrations |
| [.docs/architecture.md](.docs/architecture.md) | system architecture, runtime flow, module boundaries and dependency direction |
| [.docs/concepts/data-flow.md](.docs/concepts/data-flow.md) | request lifecycle, data ownership boundaries and cross-module flow |
| [.docs/deployment.md](.docs/deployment.md) | environment requirements, build commands, configuration policy and release checks |
| [.docs/modules/beacon-web.md](.docs/modules/beacon-web.md) | beacon-web module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/modules/beacon.md](.docs/modules/beacon.md) | beacon module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/modules/collia.md](.docs/modules/collia.md) | collia module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/modules/common.md](.docs/modules/common.md) | common module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/modules/website-web-.output-server.md](.docs/modules/website-web-.output-server.md) | website-web-.output-server module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/modules/website-web.md](.docs/modules/website-web.md) | website-web module responsibilities, implementation signals, exported-symbol hints, dependencies, state and validation |
| [.docs/project-analysis.md](.docs/project-analysis.md) | workspace inventory, documentation health, coverage model and update guidance |

## 文档健康

- 状态：pass（通过）
- 检查目标：27
- 结构失败：0
- 未纳管信号：0
- 引用错误：0
- 建议：当前无必须处理的文档治理动作；行为变更后继续运行 `/doc update`。

## 关键目录

| Directory | Inferred Role |
|-----------|---------------|
| `.docs/` | implementation documentation |
| `.github/` | supporting project directory |
| `.plans/` | SDD implementation plans |
| `.playwright-mcp/` | supporting project directory |
| `.pnpm-store/` | supporting project directory |
| `.specs/` | SDD task, status, and domain specs |
| `beacon/` | HTTP/API service module: Web/API 接入、路由注册、请求校验和服务协调 |
| `collia/` | persistence-aware service module: 数据库模型、仓储和事务边界 |
| `common/` | shared contract library: 复用 schema、数据库封装、RPC 参数/返回值和跨模块类型 |
| `deploy/` | supporting project directory |
| `website/` | frontend experience module: 页面、路由、API client、状态管理和用户交互 |

## Domain Specs

- [.specs/domain/agent-lifecycle-update.md](.specs/domain/agent-lifecycle-update.md)
- [.specs/domain/monitoring-platform.md](.specs/domain/monitoring-platform.md)

## 开发命令

```bash
task beacon
task beacon-web
task collia
task default
task website
task website-server
task website-web
cd beacon/web && npm run build
cd beacon/web && npm run dev
cd beacon/web && npm run lint
cd beacon/web && npm run preview
cd beacon/web && npm run test
cd beacon/web && npm run test:coverage
cd beacon/web && npm run test:run
cd beacon/web && npm run ts
cd website/web && npm run build
cd website/web && npm run dev
cd website/web && npm run lint
cd website/web && npm run preview
cd website/web && npm run start
cd website/web && npm run test
cd website/web && npm run test:watch
npm run build
npm run lint || true  # lint errors are non-blocking pending cleanup
npm run ts
```

## AI Agent 工作流（Spec-Driven Development）

- 新功能或复杂 Bug 先建立 `.specs/tasks/<task-id>.md` 与 `.plans/<task-id>.md`。
- Domain Spec 是最高约束来源，Plan 只描述实施路径。
- 不手工编辑 `.specs/status/*.json`，通过 SDD slash commands 维护状态。
- 老旧项目首次接入 SDD 时，可用 `/doc init --full` 与 `/doc update` 建立 `.docs`、AGENTS/CLAUDE 与 Domain Spec 初稿骨架；首次之后，`.specs/domain/` 以人工维护为准，`/doc update` 主要用于刷新 `.docs`、入口文档与索引同步。

## 技术栈

- Docker/host operation boundary
- Frontend workspace detected from positive framework/build-tool signals
- Fiber/HTTP API
- GORM persistence
- Pinia state management
- TypeScript
- Vite build tooling
- Vue frontend
- WebSocket realtime channel
- Go modules detected from `go.mod`/`go.work`
- Node.js package metadata and scripts detected from `package.json`
- Taskfile command entrypoints
- Markdown documentation under `.docs/`
- SDD domain specs under `.specs/domain/`

## 关键约束

- `.specs/domain/` 描述长期约束空间，只写必须满足的行为、不变量、状态和错误语义，不写实现方案。
- `.docs/` 描述当前实现事实，结构和内容必须能从源码、配置、构建脚本或现有文档追溯。
- `/doc update` 只全量重写带生成声明的 `.docs` 投影视图；没有生成声明的人工文档必须基于现有正文与实现事实综合更新，不得降级为通用模板。
- `AGENTS.md` 是 AI 协作入口的 SSOT；`CLAUDE.md` 是同步副本，内容必须一致。
- 新增模块、命令、路由、配置或领域概念后，重新运行 `/doc update` 并检查 `.docs/MANIFEST.yml`。
- 不在文档中写入真实密钥、Token、内部凭据或不可公开的环境值。
