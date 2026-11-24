# 架构设计

## 架构概览

`Beacon` 是包含 Go 后端和控制台前端包的应用，后端提供运行时能力，控制台通过 API 与后端协作。

当前检测到 6 个 Go workspace 模块（`beacon`、`collia`、`common`、`common/docker`、`common/logger`、`common/timex`），以及 1 个 package.json 驱动的前端包（`beacon/web/package.json`）；模块间边界以 `go.work`、各模块 `go.mod`、`package.json` 和 `.specs/domain/*.md` 为主要约束来源。

## 核心运行链路

Beacon 采用 Server-Agent 架构，存在四条明确分离的运行路径：

1. **监控查询**（Web → Server）：`beacon-web` 或第三方客户端通过 HTTP API 发起查询；Server 解析请求中的 Agent ID，从本地监控表读取对应 Agent 上报的指标，返回 JSON 响应。
2. **监控上报**（Agent → Server）：`collia` Agent 按周期采集主机与 Docker 指标，**主通道**为控制连接上的数据流——与控制流同一 gRPC 连接（Unix Socket）、独立流，批次按原子语义写入 Server 监控库；**兼容通道**为 HTTPS `POST /api/v1/host/report`（已弃用，保留若干版本），仅在数据流未就绪或推送失败时回落使用。数据流以流注册身份为准，载荷声明的 Agent 标识与流身份不一致的批次被拒绝。
3. **控制调用**（Server → Agent）：`collia` Agent 启动后主动连接 Server 在共享数据目录内监听的 Unix Socket 建立反向 gRPC tunnel（单机部署唯一承载，无 TCP 控制端口、无 mTLS）；Server 接收用户控制请求后，按 Agent ID 通过该 tunnel 调用目标 Agent 暴露的 RPC 方法（如容器操作、文件操作）。
4. **Agent 版本收敛**（Server → Agent）：Collia 注册帧携带编译版本和运行平台；Server 在注册持久化完成后异步比较镜像内置目标版本，对旧版本计算对应发布文件 SHA256 并经同一 tunnel 调用 `UpgradeAgent`。相同版本静默跳过，高版本禁止降级，结果写入系统审计。

Agent 列表中的 `online/offline` 只描述第 3 条控制链路：Server 启动会清除上一进程遗留的连接状态，并以当前进程实时持有的 tunnel 连接投影响应。第 2 条 HTTPS 监控上报的新鲜度独立判断，不能证明控制链路在线。

补充链路：

- **认证授权**：浏览器先获取 CSRF Cookie，再通过 `POST /api/v1/auth/login` 登录；JWT 只保存在 HttpOnly、SameSite Cookie 中，状态变更请求执行双提交 CSRF 校验，后续由 `UserAuthMiddleware` 与 `CasbinMiddleware` 完成身份和权限检查。新库管理员由部署期强密码创建并强制首次改密。
- **审计**：登录与认证后的变更类请求（HTTP `POST`）均写入 `s_audit` 表，可通过 `GET /api/v1/audit/query` 回查。认证后请求记录的 `Operate` 是路由名称（如「启动容器」「删除容器」），并携带操作对象与 Agent 标识（`RecordAudit` 四参写入 Username/AgentID/Operate/Object）；Collia 升级与告警等系统审计同样携带 Agent 标识。审计表由清理任务按 `AuditDays`（默认 90 天）独立保留期清理。
- **Collia 更新协调**：生产默认自动收敛、开发默认关闭；Server 对单 Agent 去重、限制全局并发并对失败自动重试设置冷却。管理员可通过专属 API 手动重试，但不能绕过互斥、架构白名单和完整性校验。
- **告警通知**：CPU、内存、磁盘统一以 `0–100` 百分数和 Server 接收时间判断完整连续窗口；未知/陈旧数据不触发也不恢复。首次告警与一次恢复先写审计，再独立尝试邮件与可选钉钉群机器人通知；持续异常每 10 分钟最多重复一次，恢复后可再次触发。邮件主题与钉钉首行共享 Server 单例消息标题。
- **WebSocket 实时通道**：`/ws/terminal`（保留的 Agent PTY 终端通道）、`/ws/:id`（容器日志流）由 `gofiber/contrib/websocket` 提供同源 upgrade；两者只使用 HttpOnly access Cookie 鉴权，不接受 query token。管理控制台当前不注册终端页面、导航或主机监控入口。

## 技术栈

| 层级 | 技术与版本 |
|------|-----------|
| 后端语言 | Go 1.25（各模块 `go.mod` 声明） |
| Web 框架 | gofiber/fiber/v2 v2.52.5 |
| WebSocket | gofiber/contrib/websocket v1.3.2 |
| ORM | gorm.io/gorm v1.25.12 + glebarez/sqlite v1.11.0 |
| 数据库驱动 | gorm.io/driver/{mysql, postgres, clickhouse} + glebarez/sqlite |
| 鉴权 | golang-jwt/jwt/v5 v5.0.0 + casbin/casbin/v2 v2.98.0 + casbin/gorm-adapter/v3 v3.28.0 |
| 配置 | spf13/viper v1.19.0 + pkg/errors v0.9.1 |
| 反向 RPC | google.golang.org/grpc v1.81.1 + google.golang.org/protobuf v1.36.11 |
| 邮件 | gomail.v2 v2.0.0 |
| 钉钉通知 | 标准库 net/http + HMAC-SHA256 签名，官方 Webhook allowlist |
| CLI | urfave/cli/v2 v2.27.4 |
| 注入 | google/wire v0.6.0 |
| Agent 端 | docker/docker v27.2.1 + shirou/gopsutil/v3 v3.24.5 + creack/pty v1.1.24 + takama/daemon v1.0.0 |
| 日志 | go.uber.org/zap v1.28.0 + lestrrat-go/file-rotatelogs v2.4.0 |
| 缓存 | patrickmn/go-cache v2.1.0 |
| 校验 | go-playground/validator/v10 v10.22.0 |
| 前端 (beacon-web) | Vue 3 + Vite 5 + TypeScript + Pinia + Element Plus + ECharts + axios + xterm.js + codemirror + UnoCSS |
| 构建/任务 | Taskfile v3（`Taskfile.yml` 入口）、pnpm 11.9 |
| 部署 | Docker + Supervisor + Nginx 反向代理 |
| 文档 | SDD 规范：`.docs/`（实现事实）+ `.specs/domain/`（长期约束） |

## 核心模块

| Module | Role |
|--------|------|
| `beacon` | **Server 端**：HTTP/API service module — Fiber 路由注册、JWT/Casbin 鉴权、监控数据持久化、Agent 生命周期编排、邮件/钉钉告警通知、反向 gRPC tunnel 客户端。承载用户可见的 API 与 WebSocket 入口。 |
| `collia` | **Agent 端**：monitoring agent module — 主机/CPU/内存/磁盘/网络指标采集（gopsutil）、Docker 操作执行（依赖 `common/docker`）、反向 gRPC tunnel 服务端、定时上报到 Server。 |
| `common` | shared contract library — 跨 Server/Agent 复用的数据库选项（`database/`）、RPC schema 与 Unix Socket 隧道协议（`rpc/`、`rpc/tunnel/`）和错误模型。 |
| `common/docker` | shared library — Docker Engine SDK 的薄封装：`IManager` 接口 + `Manager` 实现 + `ContainerSummary`/`ImageSummary`/`NetworkSummary` 摘要结构，供 collia 与 beacon 复用。 |
| `common/logger` | shared library — 基于 `go.uber.org/zap` 的日志抽象，支持文件输出、轮转、级别与格式配置；`Debug/Info/Warn/Error` 全局函数 + `Option` 函数式选项。 |
| `common/timex` | shared library — 时间工具集：`CurrentTime/CurrentDate/CurrentTimestamp/Int64ToTime` 与 `Ticker` 抽象（封装标准库 ticker 行为差异）。 |
| `beacon-web` | **Server 管理端**：Vue 3 + Vite SPA，承载 Beacon Server 的可视化控制台；包含 Pinia store、axios API client、动态路由、Codemirror、UnoCSS、ECharts 监控图表。Terminal（xterm.js）组件与依赖保留，但控制台不注册终端页面。 |

## 职责边界

| 边界 | 负责内容 | 不负责内容 | 约束来源 |
|------|----------|------------|----------|
| `beacon` | HTTP API 注册、JWT/Casbin 鉴权、监控数据持久化、审计、Agent 生命周期、反向 tunnel 客户端、WebSocket 终端与日志流 | Agent 端采集逻辑、Docker 底层操作细节、tunnel 协议内部 | `.specs/domain/monitoring-platform.md`、`go.mod`、模块 AGENTS |
| `collia` | 主机/CPU/内存/磁盘/网络指标采集、Docker 操作执行、反向 gRPC tunnel 服务端、本地配置与服务注册 | 业务持久化（仅本地预检 DB）、用户权限管理、HTTP API 注册 | `.specs/domain/agent-lifecycle-update.md`、`go.mod`、模块 AGENTS |
| `common` | 跨 Server/Agent 复用的 schema、数据库选项、RPC 参数/返回值、tunnel 协议、错误模型 | 业务逻辑、UI、HTTP 路由 | `go.mod`、导出符号、模块 AGENTS |
| `common/docker` | Docker Engine API 的薄封装：`Manager` 结构、`IManager` 接口、容器/镜像/网络/端口摘要类型 | 业务持久化、跨模块协调、HTTP 入口 | 独立 go module `github.com/amuluze/docker` |
| `common/logger` | 基于 zap 的日志抽象：级别/格式/输出/文件轮转配置、`Debug/Info/Warn/Error` 顶层 API | 业务事件追踪、链路追踪、APM 集成 | 独立 go module `github.com/amuluze/amutool/logger` |
| `common/timex` | 时间格式化与 `Ticker` 抽象：避免各模块重复实现时间戳、日期与定时器行为差异 | 业务定时任务调度、cron 表达式 | 独立 go module `github.com/amuluze/amutool/timex` |
| `beacon-web` | Beacon Server 控制台：Vue 3 SPA、Pinia 状态、API client、动态路由、Codemirror 编辑器、ECharts 监控可视化（Terminal 组件保留但未注册终端页面） | Server 端业务实现、Agent 端采集逻辑、容器运行时操作 | `beacon/web/package.json` |

## 目录结构

| Directory | Inferred Role |
|-----------|---------------|
| `.docs/` | implementation documentation |
| `.specs/` | SDD task and domain specs |
| `beacon/` | Beacon Server：HTTP/API service module — Web/API 接入、路由注册、请求校验和服务协调 |
| `collia/` | Collia Agent：monitoring agent module — 主机与 Docker 指标采集、反向 tunnel 服务端、周期上报 |
| `common/` | shared contract library — 跨 Server/Agent 复用的 schema、数据库封装、RPC 参数/返回值、tunnel 协议 |

## 调用关系

- `beacon`（Server）：导出符号信号 `func AllowPathPrefixSkipper (beacon/service/middleware/skip.go)`, `func BcryptHash (beacon/pkg/utils/hash/hash.go)`, `func BcryptVerify (beacon/pkg/utils/hash/hash.go)`；入口信号 `api := app.Group("/api")` at `beacon/service/router.go:81`、`app.Get("/health", a.healthProbe.Liveness)` at `beacon/service/router.go:87`、`app.Get("/ready", a.healthProbe.Readiness)` at `beacon/service/router.go:88`、钉钉配置路由 `registerDingTalkRoutes` at `beacon/service/router.go:371`；包含配置与外部通知边界。
- `collia`（Agent）：导出符号信号 `func BuildInjector (collia/service/wire.go)`, `func BuildInjector (collia/service/wire_gen.go)`, `func CleanupBackup (collia/service/rpc/lifecycle.go)`, `func CopyFile (collia/pkg/utils/file.go)`, `func DefaultRestartFn (collia/service/rpc/lifecycle.go)`；包含配置边界。Agent 通过 `collia/service/rpc.go` 与 `collia/service/task.go` 暴露采集与控制 RPC。
- `common`：导出符号信号 `func ComputeFreshness (common/rpc/schema/freshness.go)`, `func CurrentDate (common/timex/timex.go)`, `func CurrentTime (common/timex/timex.go)`；包含配置边界。
- `common/docker`（独立 go module `github.com/amuluze/docker`）：导出符号信号 `func NewManager (common/docker/manager.go)`, `interface IManager (common/docker/manager.go)`, `struct ContainerSummary (common/docker/container.go)`, `struct ImageSummary (common/docker/image.go)`, `struct Manager (common/docker/manager.go)`；包含配置边界。
- `common/logger`（独立 go module `github.com/amuluze/amutool/logger`）：导出符号信号 `func Debug (common/logger/default.go)`, `func DebugFields (common/logger/default.go)`, `func Debugf (common/logger/default.go)`, `func Default (common/logger/default.go)`, `func Error (common/logger/default.go)`；包含配置边界。
- `common/timex`（独立 go module `github.com/amuluze/amutool/timex`）：导出符号信号 `func CurrentDate (common/timex/timex.go)`, `func CurrentTime (common/timex/timex.go)`, `func CurrentTimestamp (common/timex/timex.go)`, `func Int64ToTime (common/timex/timex.go)`, `func NewTicker (common/timex/ticker.go)`。

## 依赖方向


### Go Module Dependency Direction

- `beacon` -> `common`（compile-time Go import）
- `beacon` -> `common/logger`（compile-time Go import）
- `beacon` -> `common/timex`（compile-time Go import）
- `collia` -> `common`（compile-time Go import）
- `collia` -> `common/docker`（compile-time Go import）
- `collia` -> `common/logger`（compile-time Go import）
- `collia` -> `common/timex`（compile-time Go import）
- `common` -> no internal Go module dependencies.
- `common/docker` -> no internal Go module dependencies.
- `common/logger` -> no internal Go module dependencies.
- `common/timex` -> no internal Go module dependencies.

### Runtime Frontend Edges

- `beacon-web` -> backend API/WebSocket boundary（runtime HTTP/WS edge, not a Go compile-time dependency）

### Go Third-Party Dependencies

- `github.com/amuluze/amutool/logger` `v0.0.0-20240821104128-caed9cc0d402`（beacon/go.mod）
- `github.com/amuluze/amutool/timex` `v0.0.0-20250508153823-fe9a5de55958`（beacon/go.mod）
- `github.com/amuluze/docker` `v0.0.0-20240822095446-429928f7463e`（collia/go.mod）
- `github.com/casbin/casbin/v2` `v2.98.0`（beacon/go.mod）
- `github.com/casbin/gorm-adapter/v3` `v3.28.0`（beacon/go.mod）
- `github.com/creack/pty` `v1.1.24`（collia/go.mod）
- `github.com/docker/docker` `v27.2.1+incompatible`（collia/go.mod）
- `github.com/docker/go-connections` `v0.5.0`（common/docker/go.mod）
- `github.com/docker/libcompose` `v0.4.0`（common/docker/go.mod）
- `github.com/glebarez/sqlite` `v1.11.0`（common/go.mod）
- `github.com/go-playground/validator/v10` `v10.22.0`（beacon/go.mod）
- `github.com/gofiber/contrib/websocket` `v1.3.2`（beacon/go.mod）
- `github.com/gofiber/fiber/v2` `v2.52.5`（beacon/go.mod）
- `github.com/golang-jwt/jwt/v5` `v5.0.0`（beacon/go.mod）
- `github.com/google/uuid` `v1.6.0`（beacon/go.mod）
- `github.com/google/wire` `v0.6.0`（beacon/go.mod）
- `github.com/lestrrat-go/file-rotatelogs` `v2.4.0+incompatible`（common/logger/go.mod）
- `github.com/patrickmn/go-cache` `v2.1.0+incompatible`（beacon/go.mod）
- `github.com/pkg/errors` `v0.9.1`（common/go.mod）
- `github.com/shirou/gopsutil/v3` `v3.24.5`（collia/go.mod）
- `github.com/spf13/viper` `v1.19.0`（beacon/go.mod）
- `github.com/stretchr/testify` `v1.9.0`（beacon/go.mod）
- `github.com/takama/daemon` `v1.0.0`（collia/go.mod）
- `github.com/tidwall/gjson` `v1.17.1`（common/docker/go.mod）
- `github.com/urfave/cli/v2` `v2.27.4`（beacon/go.mod）
- `go.uber.org/zap` `v1.28.0`（common/logger/go.mod）
- `golang.org/x/crypto` `v0.54.0`（beacon/go.mod）
- `google.golang.org/grpc` `v1.81.1`（collia/go.mod）
- `google.golang.org/protobuf` `v1.36.11`（common/go.mod）
- `gopkg.in/gomail.v2` `v2.0.0-20160411212932-81ebce5c23df`（beacon/go.mod）
- `gopkg.in/yaml.v3` `v3.0.1`（common/docker/go.mod）
- `gorm.io/driver/clickhouse` `v0.6.1`（common/go.mod）

## 配置与入口信号

- Go workspace 模块：6 个；前端包：1 个；运行时配置文件：12 个（清单见 `.docs/deployment.md`「配置说明」）。
- `beacon` Server 业务路由：77 条 `/api/v1/*` 路由，外加 `/health`、`/ready` 探针与 2 个 WebSocket 端点（`/ws/terminal`、`/ws/:id`）；完整路由清单见 `.docs/api/routes.md`。
- 详细配置语义与运行方式见 `.docs/deployment.md`，本节只保留架构级计数，避免与部署文档重复。

## 验证入口

```bash
cd beacon/web && pnpm run build
cd beacon/web && pnpm run lint
cd beacon/web && pnpm run test
cd beacon/web && pnpm run test:coverage
cd beacon/web && pnpm run test:run
cd beacon/web && pnpm run ts
```

## Domain Specs

- `.specs/domain/account-administration.md`
- `.specs/domain/agent-lifecycle-update.md`
- `.specs/domain/container-runtime.md`
- `.specs/domain/monitoring-platform.md`
