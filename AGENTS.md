# Beacon

`Beacon` 项目协作入口，汇总当前 workspace 的实现与维护事实。

## 30 秒项目摘要

`Beacon` 是 Server-Agent 架构的主机与 Docker 容器监控平台，由 `beacon` Server（含 `beacon/web` 管理控制台）与 `collia` Agent 组成，**仅支持单机部署**：Server 与 Agent 同机运行，控制通道为本机 Unix Socket 上的反向 gRPC tunnel。

核心链路：
- **监控查询**：Web 客户端 → `beacon` Server → 从本地监控表读取指定 Agent 的数据。
- **监控上报**：`collia` Agent → 控制连接（Unix Socket）上的数据流 → `beacon` Server 按批次原子写入；HTTP `POST /api/v1/host/report` 为已弃用的兼容回落。
- **控制调用**：`beacon` Server → 通过 `collia` 主动建立的反向 gRPC tunnel（Unix Socket）→ 调用目标 Agent 暴露的 RPC。

## 文档地图

| 文档 | 说明 |
|------|------|
| [DESIGN.md](DESIGN.md) | 设计意图文档：视觉基调、Token 体系、布局与组件规范、状态模型、写作与验收规范 |
| [.docs/MANIFEST.yml](.docs/MANIFEST.yml) | 文档期望清单，记录 `.docs`、`.specs/domain` 与入口文件的覆盖标准 |
| [.docs/doc-system.md](.docs/doc-system.md) | 文档系统规则与门禁边界 |
| [.docs/api/routes.md](.docs/api/routes.md) | Beacon Server HTTP 路由索引：按 group 列出方法、路径、中文名与处理函数 |
| [.docs/api/openapi.yml](.docs/api/openapi.yml) | Beacon Server HTTP API OpenAPI 3.0 契约（机器可读契约，与 routes.md 同源维护） |
| [.docs/architecture.md](.docs/architecture.md) | 系统架构：模块边界、依赖方向与三条运行链路（监控查询、指标上报、反向隧道控制） |
| [.docs/concepts/data-flow.md](.docs/concepts/data-flow.md) | 跨模块数据流：请求生命周期、数据归属边界与上报批次的原子写入路径 |
| [.docs/deployment.md](.docs/deployment.md) | 部署指南：环境要求、构建命令、配置项语义与发布前校验清单 |
| [.docs/modules/beacon-web.md](.docs/modules/beacon-web.md) | 控制台前端：Vue 3 + Element Plus + Pinia，工作区路由、主题 token 与设置工作区事实 |
| [.docs/modules/beacon.md](.docs/modules/beacon.md) | Beacon Server：Fiber 路由注册、鉴权中间件、容器与 Agent 服务及审计写入事实 |
| [.docs/modules/collia.md](.docs/modules/collia.md) | Collia Agent：主机与容器指标采集、反向 gRPC tunnel 服务端与周期上报事实 |
| [.docs/modules/common-docker.md](.docs/modules/common-docker.md) | common/docker：容器与镜像操作封装、运行时交互与错误语义事实 |
| [.docs/modules/common-logger.md](.docs/modules/common-logger.md) | common/logger：结构化日志初始化、分级输出与调用方约定 |
| [.docs/modules/common-timex.md](.docs/modules/common-timex.md) | common/timex：时间解析、时区处理与区间计算工具的实现事实 |
| [.docs/modules/common.md](.docs/modules/common.md) | common 共享契约库：跨端复用的 schema、数据库封装、RPC 参数/返回值、tunnel 协议 |
| [.docs/project-analysis.md](.docs/project-analysis.md) | workspace 盘点：模块清单、覆盖模型、文档健康计数与更新指引 |
| [.specs/domain/account-administration.md](.specs/domain/account-administration.md) | 管理员用户维护、删除目标和失败反馈的长期约束 |
| [.specs/domain/agent-lifecycle-update.md](.specs/domain/agent-lifecycle-update.md) | Agent 升级并发控制、版本推进与升级窗口数据可用性的长期约束 |
| [.specs/domain/container-runtime.md](.specs/domain/container-runtime.md) | 容器运行时能力面、列表作用域、受保护容器与删除原子性的长期约束 |
| [.specs/domain/monitoring-platform.md](.specs/domain/monitoring-platform.md) | 采样上报的数据新鲜度、阈值判定、审计记录与告警去重的长期约束 |
| [README.md](README.md) | 面向使用者的中文项目说明：功能特性、系统架构、快速开始与文档索引 |
| [README.en.md](README.en.md) | 面向使用者的英文项目说明，与中文 README 内容对齐 |
| [collia/README.md](collia/README.md) | Agent 用户入口：安装、服务生命周期、文件位置与升级语义 |
| [common/docker/README.md](common/docker/README.md) | Docker 共享库用户入口：真实运行时集成测试与资源清理 |
| [beacon/web/README.md](beacon/web/README.md) | 控制台用户入口：依赖安装、构建与开发命令 |

## 文档健康

- 状态：以当前提交执行 `task doc:check` 的退出码为准
- 检查目标：34（`.docs` 14 + Domain Spec 4 + 设计意图 1 + AI 入口与索引 10 + 项目入口说明 5）
- Domain Spec：4 份，编号约束 86 条（各 Spec 前缀唯一）
- 建议：行为变更后同步维护相关文档；检查目标与章节要求以 `.docs/MANIFEST.yml` 为准。

## 关键目录

| Directory | Inferred Role |
|-----------|---------------|
| `.docs/` | implementation documentation |
| `.specs/` | SDD 规格目录：`domain/` 存长期领域约束，`tasks/` 存过程产物（不随仓库分发） |
| `scripts/` | 文档门禁与校验脚本：Ruby 校验器（证据台账、OpenAPI 契约、门禁自测）、Python 模块信号检查器与 Go 字段形状导出器 `scripts/doc-schema` |
| `beacon/` | Beacon Server：HTTP/API service module — Web/API 接入、路由注册、请求校验和服务协调 |
| `collia/` | Collia Agent：monitoring agent module — 主机与 Docker 指标采集、反向 gRPC tunnel 服务端、周期上报 |
| `common/` | shared contract library — 跨 Server/Agent 复用的 schema、数据库封装、RPC 参数/返回值、tunnel 协议 |

## 开发命令

```bash
task beacon
task beacon-web
task collia
task default
task doc:check
task doc:verify
cd beacon/web && pnpm run build
cd beacon/web && pnpm run dev
cd beacon/web && pnpm run lint
cd beacon/web && pnpm run preview
cd beacon/web && pnpm run test
cd beacon/web && pnpm run test:coverage
cd beacon/web && pnpm run test:run
cd beacon/web && pnpm run ts
```

## AI Agent 工作流（Spec-Driven Development）

- 新功能或复杂 Bug 先建立 `.specs/tasks/<task-id>.md` 与 `.plans/<task-id>.md`；两者是过程产物，不随仓库分发，任务收口后整体清理。
- Domain Spec 是最高约束来源，Plan 只描述实施路径。
- 状态追踪由 `task doc:check` 的健康计数与本文档「文档健康」小节承担，不设单独的状态文件。
- 开发维护以私有开发仓 `amuluze/beacon-pro` 为准；本仓作为开源发布仓，GitHub Release 由发布流程同步创建（测试版标记 pre-release，正式版标记 Latest），Release 说明与 `beacon-pro` 同版本对齐。发布节奏：正式版 `vX.Y.Z`，测试版 `vX.Y.Z-beta.N`；先发测试版，验证稳定后发布同核正式版。产品镜像同时推送阿里云与 Docker Hub。
- 文档变更后运行 `task doc:check`，校验 Spec 编号唯一性与结构、MANIFEST 双向覆盖、内部链接、入口文件一致性与健康计数。

## 技术栈

- 服务端：Go + Fiber（HTTP API）+ GORM（持久化）+ Casbin（RBAC）
- Agent：Go + gopsutil（主机与容器采集）+ gRPC 反向隧道
- 控制台：Vue 3 + TypeScript + Element Plus + Pinia + Vite
- 传输：HTTP/2、WebSocket、gRPC
- 持久化：SQLite（默认）/ MySQL / PostgreSQL
- 构建与编排：Task + pnpm + Docker
- 文档：Markdown 实现文档 under `.docs/`，长期领域约束 under `.specs/domain/`

## 关键约束

- `.specs/domain/` 描述长期约束空间，只写必须满足的行为、不变量、状态和错误语义，不写实现方案。
- `.docs/` 描述当前实现事实，结构和内容必须能从源码、配置、构建脚本或现有文档追溯。
- `.docs/` 由 LLM 生成与人工维护共同产出：生成负责从源码、配置与脚本提炼结构与索引初稿，人工负责事实校正、判断性补充与边界确认；两者冲突时以可追溯的源码、配置与测试证据为准，不得用通用模板覆盖已有正文。
- `AGENTS.md` 是 AI 协作入口的 SSOT；`CLAUDE.md` 是同步副本，内容必须一致。
- 新增模块、命令、路由、配置或领域概念后，同步刷新相关项目文档并检查 `.docs/MANIFEST.yml`。
- `.specs/tasks/` 与 `.plans/` 是过程产物，不要求长久保留：任务完成后可整体清理，长期语义只沉淀在 `.specs/domain/`。
- 不在文档中写入真实密钥、Token、内部凭据或不可公开的环境值。
